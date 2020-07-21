<template>
  <container fluid>
    <div class="Tool" @mousedown="preventDefault">
      <b-navbar type="dark" variant="dark">
        <b-navbar-brand href="https://www.smartvisionworks.com/"></b-navbar-brand>
        <b-navbar-nav>
          <b-button class="nav-bar-button" v-b-toggle.sidebar-2>Instruction</b-button>
            <b-button 
              class="nav-bar-button"          
              v-b-toggle.sidebar-1
            >
              View Files
            </b-button>
        </b-navbar-nav>
      </b-navbar>
      <b-sidebar class="col-sm-2 col-md-1 sidebar" id="sidebar-1" shadow>
        <template v-slot:title>
          <div style="text-align: right" >
            Files
            <b-button @click="getFiles">
              <svg 
              width="1em" 
              height="1em" 
              viewBox="0 0 16 16" 
              class="bi bi-arrow-clockwise" 
              fill="currentColor" 
              xmlns="http://www.w3.org/2000/svg"
              >
                <path fill-rule="evenodd" d="M3.17 6.706a5 5 0 0 1 7.103-3.16.5.5 0 1 0 .454-.892A6 6 0 1 0 13.455 5.5a.5.5 0 0 0-.91.417 5 5 0 1 1-9.375.789z"/>
                <path fill-rule="evenodd" d="M8.147.146a.5.5 0 0 1 .707 0l2.5 2.5a.5.5 0 0 1 0 .708l-2.5 2.5a.5.5 0 1 1-.707-.708L10.293 3 8.147.854a.5.5 0 0 1 0-.708z"/>
              </svg>
            </b-button>
          </div>
        </template>
        <div class="px-3 py-2">
          <b-table 
            striped 
            hover 
            selectable
            select-mode="single"
            @row-selected="showFiles"
            sticky-header="500px"
            :items="folders"
            :fields="fields"
          >
          </b-table>
          <br/>
          <b-pagination
            v-model="currentPage"
            :total-rows="rows"
            :per-page="perPage"
            aria-controls="files"
          ></b-pagination>
          <p class="mt-3">Current Page: {{ currentPage }}</p>
          <b-table 
            id="files"
            striped 
            hover 
            selectable
            select-mode="single"
            @row-selected="toggle"
            sticky-header="500px"
            :per-page="perPage"
            :current-page="currentPage"
            :items="files"
            :fields="fields2"
            small
          >
            <template v-slot:cell(image)="data">
              <img :src="`/images/${activeFolder}/${data.item.image}`" style="height: 55px">
            </template>
            <template v-slot:cell(mask)="data">
              <img 
                v-if="mask_Data[activeFolder] && mask_Data[activeFolder].includes(data.item.image)"
                :src="`/masks/${activeFolder}/${data.item.image}`" 
                style="height: 55px"
              >
            </template>
          </b-table>
        </div>
      </b-sidebar>
      <b-sidebar class="sidebar" id="sidebar-2" title="Instructions" shadow>
        <div class="px-3 py-2">
          <ol>
            <li><b>1.</b> Upload Files by clicking "Upload Files" button.</li>
            <li><b>2.</b> Click and drag on the image to make a rectangle around desired object.</li>
            <li><b>3.</b> Click "Select" or press the spacebar to segment the desired area.</li>
            <li><b>4.</b> If the outcome is acceptable click "Save".</li>
            <li><b>5.</b> If adjustments need to be made use the "Foreground" and "Background" buttons or press 'F' and 'B' and draw on the image using the respective tools for foreground and background areas.</li>
            <li><b>6.</b> Click "Select" again to re-segment on the new selections.</li>
            <li><b>7.</b> Continue this process on the image until desired segmentation has been achieved.</li>
            <li><b>8.</b> Click "Save" and repeat this process for each image.</li>
            <li><b>9.</b> When all images have been completed and saved click the download button on the top toolbar to retreave your masks.</li>
            <li><b>10.</b> Good Job!</li>
          </ol>
        </div>
      </b-sidebar>
      <div style="padding: 20px">
        <canvas id="canvasOutput" ref="canvasOutput" style="width: 512px; height: 512px;" :style="{cursor: cursorType}"></canvas>
        <canvas id="canvasInput" style="width: 512px; height: 512px;"></canvas>
        <canvas id="canvasMask" style="width: 512px; height: 512px; display: none;"></canvas>
      </div>
      <b-row class="justify-content-md-center">
        <b-button-toolbar v-if="this.toolActive" key-nav aria-label="Toolbar with button groups">
          <div>
            <b-button-group class="mx-1">
              <b-button b-button v-b-tooltip.hover.bottom="'(left-arrow)'" @click="prev">
                <b-icon icon="arrow-left" aria-hidden="true"></b-icon>
              </b-button>
              <b-button b-button v-b-tooltip.hover.bottom="'(right-arrow)'"  @click="next">
                <b-icon icon="arrow-right" aria-hidden="true"></b-icon>
              </b-button>
            </b-button-group>
          </div>
          <div>
            <b-button-group  class="mx-1">
              <b-button v-b-tooltip.hover.bpointer.bottom="'(F)'" @click="fgDraw">Foreground</b-button>
              <b-button v-b-tooltip.hover.bottom="'(B)'" @click="bgDraw">Background</b-button>
              <b-button v-b-tooltip.hover.bottom="'(C)'" @click="continueDraw">Continue</b-button>
              <b-button v-b-tooltip.hover.bottom="'(U)'" @click="undo">
                <b-icon icon="arrow-counterclockwise" aria-hidden="true"></b-icon>
                Undo
              </b-button>
              <b-button v-b-tooltip.hover.bottom="'(R)'" @click="resetImg">Reset</b-button>
            </b-button-group>
          </div>
          <div>
            <b-button-group  class="mx-1">
              <b-button v-b-tooltip.hover.bottom="'(S)'" v-if="selected" @click="saveMask">Save</b-button>
            </b-button-group>
          </div>
        </b-button-toolbar>
      </b-row>
    </div>
  </container>
</template>

<script>
import paper from 'paper/dist/paper-core.min.js'
import * as cv from 'opencv.js'
import axios from 'axios'

export default {
  name: 'Tool',
  data () {
    return {
      activeFolder: '',
      perPage: 4,
      currentPage: 1,
      index: '',
      img: new Image(),
      width: 512,
      height: 512,
      mat: '',
      imageView: '',
      image: '',
      tempMask2: '',
      finalMaskResult: '',
      greenMask: '',
      redMask: '',
      addWeightedMat: '',
      img_dir: [],
      img_dir_pos: 0,
      rect: '',
      p1: '',
      p3: '',
      color: '',
      drawType: '',
      alpha: -0.5,
      beta: 1,
      gamma: 1,
      count: 0,
      drawing: false,
      drawLine: false,
      selected: false,
      continue: false,
      toolActive: false,
      rectDrawn: false,
      foregroundPoints: [],
      backgroundPoints: [],
      undoPointsFg: [],
      undoPointsBg: [],
      undoPoints: [],
      undoMats: [],
      fileNames: [],
      image_Data: null,
      mask_Data: null,
      activeIndex: -1,
      cursorType: 'crosshair',
      actionVariant: 'outline-secondary',
      downloadColor: 'warning',
      fields: [
        'folder_name',
        'images',
        'masks'
      ],
      fields2: [
        'image',
        'mask'
      ],
      folders: [],
      files: []
    }
  },
  computed: {
    activeFile() {
      if (this.activeFolder == '' || this.activeIndex == -1)
        return ''
      else
        return this.image_Data[this.activeFolder][this.activeIndex]
    },

    rows() {
      return this.files.length
    }
  },
  methods: {
    preventDefault (e) {
      e.preventDefault()
    },
    toggle (items) {
      this.activeIndex = this.image_Data[this.activeFolder].indexOf(items[0].image)
      console.log(items[0].image)
      this.activeFile = this.image_Data[this.activeFolder][this.activeIndex]
      console.log(this.image_Data[this.activeFolder][this.activeIndex])
      this.img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][this.activeIndex]}`
      this.resetImg()
    },

    getFiles () {
      let promises = []
      let self = this
      this.toolActive = true
      promises.push(
        new Promise(function(resolve, reject) {
          axios.get('/imageLocation/images')
            .then((response) => {
              console.log(response.data);
              self.image_Data = response.data
              // self.folders = response.data
              resolve()
            })
            .catch((error) => {
              reject(error)
            })
        })
      )
      promises.push(
        new Promise(function(resolve, reject) {
          axios.get('/imageLocation/masks')
            .then((response) => {
              console.log(response.data);
              self.mask_Data = response.data
              resolve()
            })
            .catch((error) => {
              reject(error)
            })
        })
      )
      Promise.all(promises).then(() => {
        this.getFolders()
      }).catch((error) => {
        console.log(error)
      })
    },

    getFolders () {
      console.log('here')
      this.folders = []
      for (let folder in this.image_Data) {
        this.folders.push({
            'folder_name': folder,
            'images': this.image_Data[folder].length,
            'masks': this.mask_Data[folder] ? this.mask_Data[folder].length : 0
          })
        console.log(this.folders)
      }
    },

    showFiles (items) {
      this.files = []
      this.activeFolder = items[0].folder_name
      let self = this
      console.log(items)
      this.image_Data[items[0].folder_name].forEach(function(filename){
        self.files.push({
          'image': filename,
          'mask': filename
        })
      })
      console.log(this.files)
    },

    saveMask () {
      this.addMask(this.finalMaskResult)
      this.actionVariant = 'outline-success'
      var canvas = document.getElementById('canvasMask')
      var self = this
      canvas.toBlob(function (blob) {
        const formData = new FormData()
        formData.append('masks', blob, `${self.activeFolder}/${self.activeFile}`)
        axios.post('/upload/masks', formData)
        .then((response) => {
          console.log(response);
        }, (error) => {
          console.log(error);
        });
        self.next()
        self.getFiles()
      })
    },

    // Canvas/Image navigation and output
    showImg () {
      this.mat = cv.imread(this.img)
      cv.resize(this.mat, this.imageView, new cv.Size(512, 512), 0, 0, cv.INTER_NEAREST)
      cv.imshow('canvasOutput', this.imageView)
      if (this.mask_Data[this.activeFolder] && this.mask_Data[this.activeFolder].includes(this.activeFile)) {
        var mask = new Image()
        mask.src = `/masks/${this.activeFolder}/${this.activeFile}`
        mask.onload = () => {
          var mat = cv.imread(mask)
          cv.imshow('canvasInput', mat)
        }
      }
    },
    prev () {
      if (this.activeIndex > 0) {
        this.activeIndex -= 1
      }
      else {
        this.activeIndex = this.image_Data[this.activeFolder].length - 1
      }
      this.img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][this.activeIndex]}`
      this.resetImg()
    },
    next () {
      if (this.activeIndex === this.image_Data[this.activeFolder].length - 1) {
        this.activeIndex = 0
      }
      else {
        this.activeIndex += 1
      }
      this.resetImg()
    },

    resetImg () {
      var blank = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(52, 64, 58, 255))
      this.cursorType = 'crosshair'
      this.actionVariant = 'outline-secondary'
      this.selected = false
      this.drawLine = false
      this.drawing = false
      this.continue = false
      this.rectDrawn = false
      this.backgroundPoints = []
      this.foregroundPoints = []
      this.undoPointsFg = []
      this.undoPointsBg = []
      this.undoPoints = []
      this.undoMats = []
      delete this.rect
      this.count = 0
      this.finalMaskResult.delete()
      this.finalMaskResult = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.addWeightedMat.delete()
      this.addWeightedMat = new cv.Mat()
      this.grabCutMask.delete()
      this.grabCutMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][this.activeIndex]}`
      this.img.onload = () => {
        this.showImg()
        cv.imshow('canvasInput', blank)
      }
    },

    // Mouse Events
    mouseDown (e) {
      this.p3 = e.point
      if (!this.drawLine) {
        this.rectangle(e)
        this.rectDrawn = true
      }
      else {
        this.draw(e)
      }
    },
    mouseDrag (e) {
      if (!this.drawLine) {
        this.rectangle(e)
        this.rectDrawn = true
      }
      else {
        this.draw(e)
      }
    },
    mouseUp () {
      if (this.foregroundPoints.length > 0 || this.backgroundPoints.length > 0) {
        var lastMat = this.image.clone()
        this.undoMats.push(lastMat)
      }
      if (this.undoPointsFg.length > 0) {
        this.undoPoints.push(this.undoPointsFg)
        this.undoPointsFg = []
      }
      if (this.undoPointsBg.length > 0) {
        this.undoPoints.push(this.undoPointsBg)
        this.undoPointsBg = []
      }
    },

    // Drawing and stuff
    draw (e) {
      this.p1 = e.point
      this.drawing = true
      cv.circle(this.image, this.p1, 2, this.color, -1)
      if (this.foregroundPoints.length === 0 && this.backgroundPoints.length === 0) {
        var tmpImg = this.image.clone()
        console.log('original image coppied')
        this.undoMats.push(tmpImg)
      }
      if (this.drawType === 'Fore point') {
        this.foregroundPoints.push(this.p1)
        this.undoPointsFg.push(this.p1)
      }
      else if (this.drawType === 'Back point') {
        this.backgroundPoints.push(this.p1)
        this.undoPointsBg.push(this.p1)
      }
      cv.imshow('canvasOutput', this.image)
    },
    fgDraw () {
      if (this.selected) {
        this.drawLine = true
        this.drawType = 'Fore point'
        this.cursorType = 'pointer'
        this.color = new cv.Scalar(255, 0, 0, 255)
      }
    },
    bgDraw () {
      if (this.selected) {
        this.drawLine = true
        this.drawType = 'Back point'
        this.cursorType = 'pointer'
        this.color = new cv.Scalar(0, 255, 0, 255)
      }
    },
    rectangle (e) {
      this.color = new cv.Scalar(0, 0, 255, 255)
      var p2 = e.point
      delete this.rect
      this.image.delete()
      if (p2.x > this.width) {
        p2.x = this.width - 2
      }
      if (p2.y > this.height) {
        p2.y = this.height - 2
      }
      p2.x = p2.x < 0 ? 1 : p2.x
      p2.y = p2.y < 0 ? 1 : p2.y
      var rectWidth = Math.abs(p2.x - this.p3.x)
      var rectHeight = Math.abs(p2.y - this.p3.y)
      this.image = this.imageView.clone()
      this.rect = new cv.Rect(Math.min(p2.x, this.p3.x), Math.min(p2.y, this.p3.y), rectWidth, rectHeight)
      cv.rectangle(this.image, p2, this.p3, this.color, 2)
      cv.imshow('canvasOutput', this.image)
    },
    continueDraw () {
      this.addMask(this.finalMaskResult)
      this.count += 1
      this.continue = true
      this.grabCutMask.delete()
      this.grabCutMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      console.log('deleted grabcutMask')
      cv.imshow('canvasMask', this.finalMaskResult)
      this.selected = false
      this.drawLine = false
      this.drawing = false
      this.foregroundPoints = []
      this.backgroundPoints = []
      this.imageView = this.addWeightedMat.clone()
      cv.imshow('canvasOutput', this.imageView)
      this.drawType = 'rect'
      this.cursorType = 'crosshair'
      this.color = new cv.Scalar(0, 0, 255, 255)
    },
    undo () {
      if (this.undoPoints.length > 0) {
        console.log('Foreground', this.foregroundPoints.length)
        console.log('Background', this.backgroundPoints.length)
        this.undoMats.splice(this.undoMats.length - 1, 1)
        cv.imshow('canvasOutput', this.undoMats[this.undoMats.length - 1])
        this.image = this.undoMats[this.undoMats.length - 1].clone()
        var points = this.undoPoints
        this.foregroundPoints = this.foregroundPoints.filter(function (el) {
          return !points[points.length - 1].includes(el)
        })
        this.backgroundPoints = this.backgroundPoints.filter(function (el) {
          return !points[points.length - 1].includes(el)
        })
        console.log('Foreground', this.foregroundPoints.length)
        console.log('Background', this.backgroundPoints.length)
        this.undoPoints.splice(this.undoPoints.length - 1, 1)
      }
    },

    // Grab-Cut and accessories
    select () {
      if (this.rectDrawn) {
        this.selected = true
        this.grab_Cut()
      }
    },
    grab_Cut () {
      var maskView = this.imageView.clone()
      var tempMat1 = new cv.Mat()
      var tempMat3 = new cv.Mat()
      var tempMat4 = new cv.Mat()
      var greenOverlay = new cv.Mat()
      var bgdModel = new cv.Mat()
      var fgdModel = new cv.Mat()
      var resizeRect = new cv.Rect(this.rect.x * 0.5, this.rect.y * 0.5, this.rect.width * 0.5, this.rect.height * 0.5)
      cv.cvtColor(maskView, maskView, 1, 0)
      var mode
      if (this.drawing) {
        console.log('foreground points', this.foregroundPoints.length)
        console.log('background points', this.backgroundPoints.length)
        mode = cv.GC_INIT_WITH_MASK
        for (var i = 0; i < this.foregroundPoints.length; i++) {
          this.tempMask2.ucharPtr(this.foregroundPoints[i].y * 0.5, this.foregroundPoints[i].x * 0.5)[0] = 1
        }
        for (i = 0; i < this.backgroundPoints.length; i++) {
          this.tempMask2.ucharPtr(this.backgroundPoints[i].y * 0.5, this.backgroundPoints[i].x * 0.5)[0] = 0
        }
      }
      else {
        mode = cv.GC_INIT_WITH_RECT
      }
      cv.resize(maskView, tempMat1, new cv.Size(256, 256), 0, 0, cv.INTER_NEAREST)
      cv.grabCut(tempMat1, this.tempMask2, resizeRect, bgdModel, fgdModel, 4, mode)
      cv.resize(this.tempMask2, tempMat3, new cv.Size(512, 512), 0, 0, cv.INTER_NEAREST)
      tempMat4 = this.createMask(tempMat3).clone()
      cv.cvtColor(tempMat4, tempMat4, 1, 0)
      cv.cvtColor(this.imageView, this.imageView, 1, 0)
      cv.addWeighted(this.greenMask, this.alpha, tempMat4.clone(), this.beta, this.gamma, greenOverlay)
      cv.addWeighted(greenOverlay.clone(), this.alpha, this.imageView, this.beta, this.gamma, this.addWeightedMat)
      cv.imshow('canvasInput', this.grabCutMask)
      cv.imshow('canvasOutput', this.addWeightedMat)
      maskView.delete(); tempMat1.delete(); tempMat3.delete(); tempMat4.delete(); greenOverlay.delete(); bgdModel.delete(); fgdModel.delete()
    },
    createMask (maskTmp) {
      var blank = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      var oneMat = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(1))
      var threeMat = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(3))
      var subMat = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(this.count))
      var foreMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0))
      cv.compare(maskTmp, oneMat, oneMat, cv.CMP_EQ)
      cv.compare(maskTmp, threeMat, threeMat, cv.CMP_EQ)
      cv.bitwise_or(oneMat, threeMat, foreMask)
      cv.subtract(foreMask, subMat, foreMask)
      cv.add(blank, foreMask, maskTmp)
      oneMat.delete(); threeMat.delete(); subMat.delete(); foreMask.delete()
      this.grabCutMask = maskTmp.clone()
      return maskTmp
    },
    addMask (finalMask) {
      cv.add(finalMask, this.grabCutMask, finalMask)
      cv.imshow('canvasMask', finalMask)
      return finalMask
    }
  },

  mounted () {
    cv['onRuntimeInitialized'] = () => {
      this.mat = new cv.Mat()
      this.imageView = new cv.Mat()
      this.image = new cv.Mat()
      this.tempMask2 = new cv.Mat()
      this.grabCutMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.finalMaskResult = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.greenMask = new cv.Mat(512, 512, cv.CV_8UC3, new cv.Scalar(152, 231, 153, 100))
      this.redMask = new cv.Mat(512, 512, cv.CV_8UC3, new cv.Scalar(231, 152, 165, 100))
      this.addWeightedMat = new cv.Mat()
      var tool = new paper.Tool()
      tool.onMouseDown = this.mouseDown
      tool.onMouseDrag = this.mouseDrag
      tool.onMouseUp = this.mouseUp
      paper.setup(this.$refs['canvasOutput'])
      this.img.onload = () => {
        this.showImg()
      }
      this.getFiles()

      window.addEventListener('keydown', (e) => {
        console.log('key:', e)
        // this.preventDefault(e)
        if (this.toolActive) {
          switch (e.code) {
            case 'ArrowUp':
              if (this.rectDrawn) {
                var greenOverlay = new cv.Mat()
                cv.cvtColor(this.grabCutMask, this.grabCutMask, 1, 0)
                cv.cvtColor(this.imageView, this.imageView, 1, 0)
                cv.addWeighted(this.greenMask, this.alpha, this.grabCutMask, this.beta, this.gamma, greenOverlay)
                cv.addWeighted(greenOverlay, this.alpha, this.imageView, this.beta, this.gamma, this.addWeightedMat)
                cv.imshow('canvasInput', this.addWeightedMat)
                greenOverlay.delete()
              }
              break
            case 'ArrowDown':
              if (this.rectDrawn) {
                var redOverlay = new cv.Mat()
                cv.cvtColor(this.grabCutMask, this.grabCutMask, 1, 0)
                cv.cvtColor(this.imageView, this.imageView, 1, 0)
                cv.addWeighted(this.redMask, this.alpha, this.grabCutMask, this.beta, this.gamma, redOverlay)
                cv.addWeighted(redOverlay, this.alpha, this.imageView, this.beta, this.gamma, this.addWeightedMat)
                cv.imshow('canvasInput', this.addWeightedMat)
                redOverlay.delete()
              }
              break
            case 'ArrowLeft':
              this.prev()
              this.resetImg()
              break
            case 'ArrowRight':
              this.next()
              this.resetImg()
              break
            case 'Space':
              this.select()
              break
            case 'KeyF':
              if (this.selected) {
                this.drawLine = true
                this.fgDraw()
              }
              break
            case 'KeyB':
              if (this.selected) {
                this.drawLine = true
                this.bgDraw()
              }
              break
            case 'KeyC':
              if (this.selected) {
                this.continueDraw()
              }
              break
            case 'KeyR':
              if (this.selected) {
                this.resetImg()
              }
              break
            case 'KeyS':
              if (this.selected) {
                this.saveMask()
              }
              break
            case 'KeyU':
              if (this.selected) {
                this.undo()
              }
              break
          }
        }
      }, false)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #539ecf;
}
canvas {
  background-color: #343a40;
  border: 10px !important;
  border-color: #343a40;
  margin: 2px;
}
.file-viewer {
  box-shadow: none !important;
  border: none !important;
  font-size: small !important;
}
.navbar-brand {
  width: 16.6666%;
  height: 60px;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100%;
  padding: 25px 0;
  float: left;
  transition: all 0.3s ease;
  background-image: url(../assets/name.svg);
  color: rgba(255, 255, 255, 0.8);
}
.nav-bar-button {
  background-color: transparent !important;
  border-color: transparent !important;
}
.nav-bar-button:hover {
  background-color: #5a6268 !important;
  border-color: #5a6268 !important;
}
.b-sidebar-body{
  background-color: black !important;
}
.b-sidebar-btn-secondary {
  width: 100%;
}
.card-main {
  margin-top: 20px;
  background-color: #343a40;
  border: none;
}
</style>
