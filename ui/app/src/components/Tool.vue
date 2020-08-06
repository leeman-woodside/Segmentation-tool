<template>
  <b-container fluid>
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
          <div>
            Files
            <b-button class="refresh-button" @click="getFiles">
              <b-icon-arrow-clockwise></b-icon-arrow-clockwise>
            </b-button>
          </div>
        </template>
        <div class="px-3 py-2">
          <b-table
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
              <b-icon-check2-square v-if="mask_Data[activeFolder] && mask_Data[activeFolder].includes(data.item.image)" style="color: green"></b-icon-check2-square>
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
            <li style="color:red;"><h6><b>6.</b></h6> If the outcome is acceptable click "Save".</li>
            <li><h6><b>7.</b></h6> If adjustments need to be made use the "Foreground" and "Background" buttons or press 'F' and 'B' on the keyboard and using the mouse, draw on the image using the respective tools for foreground and background areas.</li>
            <li><h6><b>8.</b></h6> Click "Select" again to re-segment on the new selections.</li>
            <li><h6><b>9.</b></h6> Continue this process on the image until desired segmentation has been achieved.</li>
            <li><h6><b>10.</b></h6> If you make a mistake using the foreground and background tools you can click "Undo" or press 'U' on the keyboard to go back to a previous selection.</li>
            <li style="color:red;"><h6><b>11.</b></h6> Click "Save" and repeat this process for each image.</li>
            <li><h6><b>12.</b></h6> Good Job!</li>
          </ol>
        </div>
      </b-sidebar>
      <div style="padding: 20px">
        <canvas id="canvasOutput" ref="canvasOutput" style="width: 512px; height: 512px;" :style="{cursor: cursorType}"></canvas>
        <canvas id="canvasInput" style="width: 512px; height: 512px;"></canvas>
        <canvas id="canvasMask" style="width: 512px; height: 512px; display: none;"></canvas>
      </div>
      <b-row class="justify-content-md-center">
        <b-button-toolbar v-if="this.activeFile != ''" key-nav aria-label="Toolbar with button groups">
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
              <b-button v-b-tooltip.hover.bottom="'(spacebar)'" @click="select">Run Tool</b-button>
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
  </b-container>
</template>

<script>
import paper from 'paper/dist/paper-core.min.js'
import * as cv from 'opencv.js-webassembly'
import axios from 'axios'

export default {
  name: 'Tool',
  data () {
    return {
      // Constants
      width: 512,
      height: 512,
      alpha: -0.3,
      beta: 1,
      gamma: 1,

      scrollY: 0,
      scrollX: 0,
      activeFolder: '',
      bufferArray: [],
      img: new Image(),
      masterMat: '',
      originalImage: '',
      updateMat: '',
      maskViewFinal: '',
      rect: '',
      rectPoint1: {},
      rectPoint2: {},
      drawColor: '',
      rectColor: '',
      drawType: '',
      maskObjCount: 0,
      drawing: false,
      drawLine: false,
      selected: false,
      rectDrawn: false,
      points: {fg:[], bg:[]},
      foregroundPoints: [],
      backgroundPoints: [],
      undoPoints: [],
      undoMats: [],
      image_Data: null,
      mask_Data: null,
      activeIndex: -1,
      cursorType: 'crosshair',
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
    //disables a keyboard key's default action
    preventDefault (e) {
      e.preventDefault()
    },
    //sets the index and loads the buffer array based on the file that is selected from the table
    toggle (items) {
      if (items.length === 0) {
        return 
      }
      this.activeIndex = this.image_Data[this.activeFolder].indexOf(items[0].image)
      this.createBufferArray()
      this.resetImg()
    },
    //this is the array of images that load ahead and behind of the current image being worked on
    createBufferArray() {
      this.bufferArray = []
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
    },
    //gets the images file paths from the server and puts that into two objects 
    getFiles () {
      let promises = []
      let self = this
      promises.push(
        new Promise(function(resolve, reject) {
          axios.get('/imageLocation/images')
            .then((response) => {
              self.image_Data = response.data
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
    //populates the folder table for the sidebar
    getFolders () {
      this.folders = []
      for (let folder in this.image_Data) {
        this.folders.push({
            'folder_name': folder,
            'images': this.image_Data[folder].length,
            'masks': this.mask_Data[folder] ? this.mask_Data[folder].length : 0
          })
      }
    },
    //populates the files table from whichever folder has been selected
    showFiles (items) {
      if (items.length === 0) {
        return
      }
      this.files = []
      this.activeFolder = items[0].folder_name
      let self = this
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
    },
    //saves a finished mask to the server
    saveMask () {
      this.maskViewFinal = this.addMask(this.maskViewFinal)
      cv.imshow('canvasMask', this.maskViewFinal)
      var canvas = document.getElementById('canvasMask')
      var self = this
      canvas.toBlob(function (blob) {
        const formData = new FormData()
        formData.append('masks', blob, `${self.activeFolder}/${self.activeFile}`)
        console.log(formData)
        axios.post('/upload/masks', formData)
        .then((response) => {
          console.log(response);
          self.getFiles()
        })
        .catch((error) => {
          console.log(error);
        });
        self.next()
      })
    },
    // Images are shown from their file path on the server
    showImg () {
      this.originalImage = cv.imread(this.img)
      cv.resize(this.originalImage, this.originalImage, new cv.Size(512, 512), 0, 0, cv.INTER_NEAREST)
      this.masterMat = this.originalImage.clone()
      cv.imshow('canvasOutput', this.masterMat)
      if (this.mask_Data[this.activeFolder] && this.mask_Data[this.activeFolder].includes(this.activeFile)) {
        var mask = new Image()
        mask.src = `/masks/${this.activeFolder}/${this.activeFile}?${Math.random()}`
        mask.onload = () => {
          var mat = cv.imread(mask)
          cv.imshow('canvasInput', mat)
        }
      }
    },
    // navagation thru the different images on the server
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
      //unshift() adding to the front and pop() popping off the back
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
      this.resetImg()
    },
    //if all goes wrong on an image you can just hit reset to start over on that image
    resetImg () {
      let blank = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(52, 64, 58, 255))
      this.cursorType = 'crosshair'
      this.selected = false
      this.drawLine = false
      this.drawing = false
      this.rectDrawn = false
      this.undoPoints = []
      this.undoMats = []
      delete this.rect
      this.maskObjCount = 0
      this.grabCutMask.delete()
      this.grabCutMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.maskViewFinal.delete()
      this.maskViewFinal = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.img.src = this.bufferArray[3].src
      this.img.onload = () => {
        this.showImg()
        cv.imshow('canvasInput', blank)
      }
    },
    // Mouse Events
    mouseDown (e) {
      if (!this.drawLine && this.activeFile != '') {
        this.rectPoint1.x = e.point.x + this.scrollX
        this.rectPoint1.y = e.point.y + this.scrollY
        console.log('setting rectPoint1: ', this.rectPoint1)
        this.rectangle(e)
        this.rectDrawn = true
      }
      if (this.drawLine) {
        if (this.points.fg.length === 0 && this.points.bg.length === 0) {
          var tmpImg = this.imageDraw.clone()
          this.undoMats.push(tmpImg)
          console.log('pushed first mat to undoMat array baby')
        }
        this.fg_bg_PointArrayTracker (e.point)
        this.draw(e)
      }
    },
    mouseDrag (e) {
      if (!this.drawLine && this.activeFile != '') {
        this.rectangle(e)
        this.rectDrawn = true
      }
      if (this.drawLine) {
        this.draw(e)
      }
    },
    mouseUp () {
      var lastMat = this.imageDraw.clone()
      if (this.foregroundPoints.length > 0) {
        this.undoMats.push(lastMat)
        this.undoPoints.push(this.foregroundPoints)
        this.points.fg.push(this.foregroundPoints)
        this.foregroundPoints = []
      }
      if (this.backgroundPoints.length > 0) {
        this.undoMats.push(lastMat)
        this.undoPoints.push(this.backgroundPoints)
        this.points.bg.push(this.backgroundPoints)
        this.backgroundPoints = []
      }
    },
    //ROI
    rectangle (e) {
      this.rectColor = new cv.Scalar(110, 170, 255, 255)
      delete this.rect
      this.imageDraw.delete()
      this.rectPoint2.x = e.point.x + this.scrollX
      this.rectPoint2.y = e.point.y + this.scrollY
      this.rectPoint2.x = this.rectPoint2.x > this.width ? this.width - 2 : this.rectPoint2.x
      this.rectPoint2.y = this.rectPoint2.y > this.height ? this.height - 2 : this.rectPoint2.y
      this.rectPoint2.x = this.rectPoint2.x < 0 ? 1 : this.rectPoint2.x
      this.rectPoint2.y = this.rectPoint2.y < 0 ? 1 : this.rectPoint2.y
      var rectWidth = Math.abs(this.rectPoint2.x - this.rectPoint1.x)
      var rectHeight = Math.abs(this.rectPoint2.y - this.rectPoint1.y)
      this.imageDraw = this.masterMat.clone()
      this.rect = new cv.Rect(Math.min(this.rectPoint2.x, this.rectPoint1.x), Math.min(this.rectPoint2.y, this.rectPoint1.y), rectWidth, rectHeight)
      cv.rectangle(this.imageDraw, this.rectPoint1, this.rectPoint2, this.rectColor, 2)
      cv.imshow('canvasOutput', this.imageDraw)
    },
    //creates arrays to populate points object
    //changes the pixel class of the updateMat
    fg_bg_PointArrayTracker (p) {
      if (this.drawLine && p.x > this.rect.x && p.y > this.rect.y && p.x < (this.rect.x + this.rect.width) && p.y < (this.rect.y + this.rect.height)) {
        if (this.drawType === 'Fore point') {
          this.foregroundPoints.push(p)
          this.updateMat.ucharPtr(p.y * 0.5, p.x * 0.5)[0] = 1
          
        }
        else if (this.drawType === 'Back point') {
          this.backgroundPoints.push(p)
          this.updateMat.ucharPtr(p.y * 0.5, p.x * 0.5)[0] = 0
        }
      }
    },
    //draws points on canvas
    draw (e) {
      let drawPoint = {}
      drawPoint.x = e.point.x + this.scrollX
      drawPoint.y = e.point.y + this.scrollY
      //allow drawing only in ROI
      if (drawPoint.x > this.rect.x && drawPoint.y > this.rect.y && drawPoint.x < (this.rect.x + this.rect.width) && drawPoint.y < (this.rect.y + this.rect.height)) {
        this.drawing = true
        cv.circle(this.imageDraw, drawPoint, 1, this.drawColor, -1)
        this.fg_bg_PointArrayTracker (drawPoint)
        cv.imshow('canvasOutput', this.imageDraw)
      }
    },
    //foreground points
    fgDraw () {
      if (this.selected) {
        this.drawLine = true
        this.drawType = 'Fore point'
        this.cursorType = 'cell'
        this.drawColor = new cv.Scalar(255, 0, 0, 255)
      }
    },
    //background points
    bgDraw () {
      if (this.selected) {
        this.drawLine = true
        this.drawType = 'Back point'
        this.cursorType = 'cell'
        this.drawColor = new cv.Scalar(0, 255, 0, 255)
      }
    },
    //if more than one object requires segmentation in an image
    continueDraw () {
      this.maskViewFinal = this.addMask(this.maskViewFinal)
      this.maskObjCount += 1
      this.selected = false
      this.drawLine = false
      this.drawing = false
      let weightedMat = new cv.Mat()
      let gcMask = this.grabCutMask.clone()
      let originalMat = this.originalImage.clone()
      let maskFinal = this.maskViewFinal.clone()
      cv.cvtColor(gcMask, gcMask, 1, 0)
      cv.cvtColor(originalMat, originalMat, 1, 0)
      cv.cvtColor(maskFinal, maskFinal, 1, 0)
      cv.addWeighted(maskFinal, 0.3, originalMat, this.beta, this.gamma, weightedMat)
      this.masterMat = weightedMat.clone()
      cv.imshow('canvasOutput', this.masterMat)
      this.drawType = 'rect'
      this.cursorType = 'crosshair'
      this.rectColor = new cv.Scalar(110, 170, 255, 255)
    },
    //takes away last drawn points up to drawing the ROI
    undo () {
      if (this.undoPoints.length > 0) {
        console.log('Foreground', this.points.fg)
        console.log('Background', this.points.bg)
        console.log(this.undoMats)
        this.undoMats.pop()
        cv.imshow('canvasOutput', this.undoMats[this.undoMats.length - 1])
        let undoPts = this.undoPoints[this.undoPoints.length - 1]
        this.imageDraw = this.undoMats[this.undoMats.length - 1].clone()
        if (this.points.fg.includes(undoPts)) {
          for (var i = 0; i < undoPts.length; i++) {
            this.updateMat.ucharPtr(undoPts[i].y * 0.5, undoPts[i].x * 0.5)[0] = 2
          }
          this.points.fg.pop()
        }
        else if (this.points.bg.includes(this.undoPoints[this.undoPoints.length - 1])) {
          for (var j = 0; j < undoPts.length; j++) {
            this.updateMat.ucharPtr(undoPts[j].y * 0.5, undoPts[j].x * 0.5)[0] = 2
          }
          this.points.bg.pop()
        }
        console.log('Foreground', this.points.fg)
        console.log('Background', this.points.bg)
        this.undoPoints.pop()
        this.select()
      }
    },
    // selecting ROI
    select () {
      if (this.rectDrawn) {
        let greenMask = new cv.Mat(512, 512, cv.CV_8UC3, new cv.Scalar(152, 231, 153, 100))
        this.rectColor = new cv.Scalar(0, 0, 255, 255)
        this.selected = true
        var weightedMat = new cv.Mat()
        var workingOverlay = new cv.Mat()
        this.grabCutMask = this.grab_Cut()
        var mask = this.grabCutMask.clone()
        var finalMasks = this.maskViewFinal.clone()
        cv.cvtColor(mask, mask, 1, 0)
        cv.cvtColor(this.originalImage, this.originalImage, 1, 0)
        cv.cvtColor(finalMasks, finalMasks, 1, 0)
        cv.addWeighted(greenMask, this.alpha, mask, this.beta, this.gamma, workingOverlay)
        cv.add(workingOverlay, finalMasks, workingOverlay)
        cv.addWeighted(workingOverlay, this.alpha, this.originalImage, this.beta, this.gamma, weightedMat)
        cv.rectangle(weightedMat, this.rectPoint1, this.rectPoint2, this.rectColor, 2)
        cv.add(finalMasks, mask, finalMasks)
        cv.imshow('canvasInput', finalMasks)
        cv.imshow('canvasOutput', weightedMat)
        cv.imshow('canvasMask', this.maskViewFinal)
        workingOverlay.delete()
      }
    },
    //perform grabcut 
    grab_Cut () {
      var maskView = this.masterMat.clone()
      var bgdModel = new cv.Mat()
      var fgdModel = new cv.Mat()
      var mode
      var resizeRect = new cv.Rect(this.rect.x * 0.5, this.rect.y * 0.5, this.rect.width * 0.5, this.rect.height * 0.5)
      cv.cvtColor(maskView, maskView, 1, 0)
      if (this.drawing) {
        mode = cv.GC_INIT_WITH_MASK
      }
      else {
        mode = cv.GC_INIT_WITH_RECT
      }
      cv.resize(maskView, maskView, new cv.Size(256, 256), 0, 0, cv.INTER_NEAREST)
      cv.grabCut(maskView, this.updateMat, resizeRect, bgdModel, fgdModel, 4, mode)
      cv.resize(this.updateMat, maskView, new cv.Size(512, 512), 0, 0, cv.INTER_NEAREST)
      maskView = this.createMask(maskView).clone()
      bgdModel.delete(); fgdModel.delete()
      return maskView
    },
    //make black and white visual mask
    createMask (maskTmp) {
      var blank = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      var isFGMat = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(cv.GC_FGD))
      var probFGMat = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(cv.GC_PR_FGD))
      var subMat = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(this.maskObjCount))
      var foreMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0))
      cv.compare(maskTmp, isFGMat, isFGMat, cv.CMP_EQ)
      cv.compare(maskTmp, probFGMat, probFGMat, cv.CMP_EQ)
      cv.bitwise_or(isFGMat, probFGMat, foreMask)
      cv.subtract(foreMask, subMat, foreMask)
      cv.add(blank, foreMask, maskTmp)
      isFGMat.delete(); probFGMat.delete(); subMat.delete(); foreMask.delete()
      return maskTmp
    },
    //adds current grabCutMask with the final mask
    addMask (mask) {
      cv.add(mask, this.grabCutMask, mask)
      cv.imshow('canvasMask', mask)
      return mask
    }
  },

  created () {
    document.body.addEventListener('scroll', (e) => {
      this.scrollX = e.target.scrollLeft
      this.scrollY = e.target.scrollTop
    });
  },

  mounted () {
    cv['onRuntimeInitialized'] = () => { 
      this.getFiles()
      this.$Progress.finish()
      this.masterMat = new cv.Mat()
      this.originalImage = new cv.Mat()
      this.updateMat = new cv.Mat()
      this.imageDraw = new cv.Mat()
      this.grabCutMask = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      this.maskViewFinal = new cv.Mat(512, 512, cv.CV_8UC1, new cv.Scalar(0, 0, 0, 255))
      var tool = new paper.Tool()
      tool.onMouseDown = this.mouseDown
      tool.onMouseDrag = this.mouseDrag
      tool.onMouseUp = this.mouseUp
      paper.setup(this.$refs['canvasOutput'])
      this.img.onload = () => {
        this.showImg()
      }

      window.addEventListener('keydown', (e) => {
        console.log('key:', e)
        if (this.activeFile != '') {
          switch (e.code) {
            // case 'ArrowUp':
            //   if (this.rectDrawn) {
            //     var workingOverlay = new cv.Mat()
            //     cv.cvtColor(this.grabCutMask, this.grabCutMask, 1, 0)
            //     cv.cvtColor(this.masterMat, this.masterMat, 1, 0)
            //     cv.addWeighted(this.greenMask, this.alpha, this.grabCutMask, this.beta, this.gamma, workingOverlay)
            //     cv.addWeighted(workingOverlay, this.alpha, this.masterMat, this.beta, this.gamma, this.addWeightedMat)
            //     cv.imshow('canvasInput', this.addWeightedMat)
            //     workingOverlay.delete()
            //   }
            //   break
            // case 'ArrowDown':
            //   if (this.rectDrawn) {
            //     var redOverlay = new cv.Mat()
            //     cv.cvtColor(this.grabCutMask, this.grabCutMask, 1, 0)
            //     cv.cvtColor(this.masterMat, this.masterMat, 1, 0)
            //     cv.addWeighted(this.redMask, this.alpha, this.grabCutMask, this.beta, this.gamma, redOverlay)
            //     cv.addWeighted(redOverlay, this.alpha, this.masterMat, this.beta, this.gamma, this.addWeightedMat)
            //     cv.imshow('canvasInput', this.addWeightedMat)
            //     redOverlay.delete()
            //   }
            //   break
            case 'ArrowLeft':
              this.prev()
              break
            case 'ArrowRight':
              this.next()
              break
            case 'Space':
              this.preventDefault(e)
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
.refresh-button {
  background: none;
  border: none; 
  color: #95999d;
}
.refresh-button:hover {
  color: #343a40 !important;
}
.no-js #loader { display: none;  }
.js #loader { display: block; position: absolute; left: 100px; top: 0; }
</style>
