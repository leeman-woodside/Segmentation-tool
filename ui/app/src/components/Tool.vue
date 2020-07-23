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
            :activeIndex = 0
            selectable
            select-mode="single"
            @row-selected="showFiles"
            sticky-header="500px"
            :items="folders"
            :fields="fields"
          >
          </b-table>
          <br/>
          <b-table 
            id="files"
            striped 
            hover 
            selectable
            select-mode="single"
            @row-selected="toggle"
            sticky-header="500px"
            :items="files"
            :fields="fields2"
            small
          >
            <template v-slot:cell(mask)="data">
              <b-icon-check2-square v-if="mask_Data[activeFolder] && mask_Data[activeFolder].includes(data.item.image)"></b-icon-check2-square>
            </template>
          </b-table>
        </div>
      </b-sidebar>
      <b-sidebar class="sidebar" id="sidebar-2" title="Instructions" shadow>
        <div class="px-3 py-2">
          <ol style="text-align: left">
            <li><h6><b>1.</b></h6> Click "View Files".</li>
            <li><h6><b>2.</b></h6> In the sidebar menu, select in the top table which folder of images you want to work on.</li>
            <li><h6><b>3.</b></h6> In the second table that opens up select which image you want to start with.</li>
            <li><h6><b>4.</b></h6> Click and drag on the image to make a rectangle around desired object.</li>
            <li><h6><b>5.</b></h6>Click "Select" or press the spacebar to segment the desired area.</li>
            <li><h6><b>6.</b></h6> If the outcome is acceptable click "Save".</li>
            <li><h6><b>7.</b></h6> If adjustments need to be made use the "Foreground" and "Background" buttons or press 'F' and 'B' on the keyboard and draw on the image using the respective tools for foreground and background areas.</li>
            <li><h6><b>8.</b></h6> Click "Select" again to re-segment on the new selections.</li>
            <li><h6><b>9.</b></h6> Continue this process on the image until desired segmentation has been achieved.</li>
            <li><h6><b>10.</b></h6> Click "Save" and repeat this process for each image.</li>
            <li><h6><b>11.</b></h6> Good Job!</li>
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
              <b-button v-b-tooltip.hover.bottom="'(spacebar)'" @click="select">Select</b-button>
              <b-button v-b-tooltip.hover.bpointer.bottom="'(F)'" @click="fgDraw">Foreground</b-button>
              <b-button v-b-tooltip.hover.bottom="'(B)'" @click="bgDraw">Background</b-button>
              <b-button v-b-tooltip.hover.bottom="'(C)'" @click="continueDraw">Add Object</b-button>
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
import * as cv from 'opencv.js-webassembly'
import axios from 'axios'

export default {
  name: 'Tool',
  data () {
    return {
      activeFolder: '',
      index: '',
      bufferArray: [],
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
      p2: '',
      drawColor: '',
      rectColor: '',
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
    }
  },
  methods: {
    preventDefault (e) {
      e.preventDefault()
    },
    toggle (items) {
      this.activeIndex = this.image_Data[this.activeFolder].indexOf(items[0].image)
      this.activeFile = this.image_Data[this.activeFolder][this.activeIndex]
      this.createBufferArray()
      this.resetImg()
    },

    createBufferArray() {
      this.bufferArray = []
      console.log('creating buffer array')
      let popIndex = this.activeIndex
      let unshiftIndex = this.activeIndex
      let img1 = new Image()
      img1.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][this.activeIndex]}`
      this.bufferArray.push(img1)
      for (var i = 1; i < 4; i++) {
        let img2 = new Image()
        if (popIndex === this.image_Data[this.activeFolder].length - 1) {
          popIndex = 0
          img2.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][popIndex]}`
        }
        else {
          popIndex = popIndex + 1
          img2.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][popIndex]}`
        }
        this.bufferArray.push(img2)
      }
      for (var j = 1; j < 4; j++) {
        let img3 = new Image()
        if (unshiftIndex === 0) {
          unshiftIndex = this.image_Data[this.activeFolder].length - 1
          img3.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][unshiftIndex]}`
        }
        else {
          unshiftIndex = unshiftIndex - 1
          img3.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][unshiftIndex]}`
        }
        this.bufferArray.unshift(img3)
      }
      this.img.src = this.bufferArray[3].src
      for (var k = 0; k < this.bufferArray.length; k ++) {
        console.log(this.bufferArray[k].src)
      }
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
        if (self.mask_Data[self.activeFolder] && self.mask_Data[self.activeFolder].includes(filename)) {
          self.files.push({
            'image': filename,
            'mask': filename
          })
        }
        else {
          self.files.push({
            'image': filename,
            'mask': ''
          })
        }
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
      this.bufferArray.pop()
      let img = new Image()
      let index = this.activeIndex
      if (this.activeIndex > 0) {
        this.activeIndex -= 1
      }
      else {
        this.activeIndex = this.image_Data[this.activeFolder].length - 1
      }
      //if backwards unshift() adding to the front and pop() popping off the back
      if (this.activeIndex - 3 < 0) {
        index = (this.image_Data[this.activeFolder].length - 1) + (this.activeIndex - 3)
        img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][index + 1]}`
        this.bufferArray.unshift(img)
      }
      else {
        index = this.activeIndex - 3
        img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][index]}`
        this.bufferArray.unshift(img)
      }
      for (var k = 0; k < this.bufferArray.length; k ++) {
        console.log(this.bufferArray[k].src)
      }
      this.resetImg()
    },
    next () {
      this.bufferArray.shift()
      let img = new Image()
      let index = this.activeIndex
      if (this.activeIndex === this.image_Data[this.activeFolder].length - 1) {
        this.activeIndex = 0
      }
      else {
        this.activeIndex += 1
      }
      //if forwards shift() off the front and push() to the back
      if (this.activeIndex + 3 > this.image_Data[this.activeFolder].length - 1) {
        index = (this.activeIndex + 3) - (this.image_Data[this.activeFolder].length - 1)
        img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][index - 1]}`
        this.bufferArray.push(img)
      }
      else {
        index = this.activeIndex + 3
        img.src = `/images/${this.activeFolder}/${this.image_Data[this.activeFolder][index]}`
        this.bufferArray.push(img)
      }
      for (var k = 0; k < this.bufferArray.length; k ++) {
        console.log(this.bufferArray[k].src)
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
      this.img.src = this.bufferArray[3].src
      this.img.onload = () => {
        this.showImg()
        cv.imshow('canvasInput', blank)
      }
    },

    // Mouse Events
    mouseDown (e) {
      if (!this.drawLine) {
        this.p3 = e.point
        this.rectangle(e)
        this.rectDrawn = true
      }
      else {
        this.draw(e)
      }
    },
    mouseDrag (e) {
      if (e.event.which == 2 || e.event.which == 4) {
        paper.view.center = e.downPoint.subtract(e.point).add(paper.view.center);
      }
      else {
        if (!this.drawLine) {
          this.rectangle(e)
          this.rectDrawn = true
        }
        else {
          this.draw(e)
        }
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
    // mouseWheel (e) {
    //   // Store previous view state.
    // var oldZoom = paper.view.zoom;
    // var oldCenter = paper.view.center;

    // // Get mouse position.
    // // It needs to be converted into project coordinates system.
    // var mousePosition = paper.view.viewToProject(new paper.Point(e.offsetX, e.offsetY));

    // // Update view zoom.
    // var newZoom = e.deltaY > 0
    //     ? oldZoom * 1.5
    //     : oldZoom / 1.5;
    // paper.view.zoom = newZoom;

    // // Update view position.
    // paper.view.center += (mousePosition - oldCenter) * (1 - (oldZoom / newZoom));
    // },

    // Drawing and stuff
    draw (e) {
      this.p1 = e.point
      this.drawing = true
      cv.circle(this.image, this.p1, 2, this.drawColor, -1)
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
        this.drawColor = new cv.Scalar(255, 0, 0, 255)
      }
    },
    bgDraw () {
      if (this.selected) {
        this.drawLine = true
        this.drawType = 'Back point'
        this.cursorType = 'pointer'
        this.drawColor = new cv.Scalar(0, 255, 0, 255)
      }
    },
    rectangle (e) {
      this.rectColor = new cv.Scalar(0, 0, 255, 255)
      this.p2 = e.point
      delete this.rect
      this.image.delete()
      if (this.p2.x > this.width) {
        this.p2.x = this.width - 2
      }
      if (this.p2.y > this.height) {
        this.p2.y = this.height - 2
      }
      this.p2.x = this.p2.x < 0 ? 1 : this.p2.x
      this.p2.y = this.p2.y < 0 ? 1 : this.p2.y
      var rectWidth = Math.abs(this.p2.x - this.p3.x)
      var rectHeight = Math.abs(this.p2.y - this.p3.y)
      this.image = this.imageView.clone()
      this.rect = new cv.Rect(Math.min(this.p2.x, this.p3.x), Math.min(this.p2.y, this.p3.y), rectWidth, rectHeight)
      cv.rectangle(this.image, this.p2, this.p3, this.rectColor, 2)
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
      // this.imageView = this.addWeightedMat.clone()
      cv.imshow('canvasOutput', this.imageView)
      this.drawType = 'rect'
      this.cursorType = 'crosshair'
      this.rectColor = new cv.Scalar(0, 0, 255, 255)
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
      cv.rectangle(this.addWeightedMat, this.p3, this.p2, this.rectColor, 2)
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
    },

    mouseEvent () {
      console.log('hi')
    }
  },

  mounted () {
    cv['onRuntimeInitialized'] = () => {
      console.log(cv)
      this.getFiles()
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
      tool.mousewheel = this.mouseWheel
      paper.setup(this.$refs['canvasOutput'])
      this.img.onload = () => {
        this.showImg()
      }

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
