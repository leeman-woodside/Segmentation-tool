<template>
<b-container fluid>
  <div class="Tool" @mousedown="preventDefault">
    <b-navbar type="dark" variant="dark">
      <b-navbar-brand href="http://localhost:8082/#/"></b-navbar-brand>
      <b-navbar-nav>
          <input
            style="display: none"
            name="myFiles"
            type="file"
            @change="onFileSelected"
            webkitdirectory
            ref="fileInput"/>
          <b-button class="nav-bar-button" @click="$refs.fileInput.click()">
            Upload Files
            <b-icon icon="folder-plus" aria-hidden="true"></b-icon>
          </b-button>
          <b-button class="nav-bar-button" style="position: right" @click="getFiles">
            Get Masks
          </b-button>
          <b-button class="nav-bar-button" style="position: right" @click="generateZip">
            <b-icon icon="download" aria-hidden="true"></b-icon>
            <b-badge :variant="downloadColor" v-if="numSaved > 0">{{ numSaved + " / " + img_dir.length }}</b-badge>
          </b-button>
      </b-navbar-nav>
    </b-navbar>
    <b-card class="card-main">
      <b-row>
        <b-col cols="6">
          <b-card
            border-variant="default"
            header-bg-variant="default"
            header-text-variant="default"
            align="left"
          >
            <div slot="header">
              <b-row>
                <b-col cols="4">
                  <h2>Folders</h2>
                </b-col>
              </b-row>
            </div>
            <b-collapse id="cardSystems" visible>
              <b-card-text>
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
              </b-card-text>
            </b-collapse>
          </b-card>
        </b-col>
        <b-col cols="6">
          <b-card
            border-variant="default"
            header-bg-variant="default"
            header-text-variant="default"
            align="left"
          >
            <div slot="header">
              <b-row>
                <b-col cols="4">
                  <h2>Files</h2>
                </b-col>
              </b-row>
            </div>
            <b-collapse id="fileSystems" visible>
              <b-card-text>
                <b-table 
                  striped 
                  hover 
                  selectable
                  select-mode="single"
                  sticky-header="500px"
                  :items="files"
                  :fields="fields2"
                >
                  <template v-slot:cell(image)="data">
                    <img :src="`/images/${activeFolder}/${data.item.image}`" style="height: 141px">
                  </template>
                  <template v-slot:cell(mask)="data">
                    <img 
                      v-if="mask_Data[activeFolder] && mask_Data[activeFolder].includes(data.item.image)"
                      :src="`/masks/${activeFolder}/${data.item.image}`" 
                      style="height: 141px"
                    >
                  </template>
                </b-table>
              </b-card-text>
            </b-collapse>
          </b-card>
        </b-col>
      </b-row>
    </b-card>
  </div>
</b-container>
</template>

<script>
import axios from 'axios'
import * as FileSaver from 'FileSaver'
import * as Zip from 'jszip'
import * as ZipUtil from 'jszip-utils'

export default {
  name: 'Admin',
  data () {
    return {
      activeFolder: '',
      mask_Data: null,
      image_Data: null,
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
  methods: {
    preventDefault (e) {
      e.preventDefault()
    },

    onFileSelected (e) {
      var files = e.target.files
      for (let i = 0; i < files.length; i++) {
        let formData = new FormData()
        formData.append('images', e.target.files[i])
        axios.post('/upload/images', formData)
        .then((response) => {
          console.log(response);
        }, (error) => {
          console.log(error);
        });
          //mongoDB stuff
        // axios.post('/image', {
          //   imgName: fName,
        //   imgType: fType,
        //   imgSize: fSize
        // })
        // .then((response) => {
          //   console.log(response);
        // }, (error) => {
          //   console.log(error);
        // });
      }
    },

    getFiles () {
      let promises = []
      let self = this
      promises.push(
        new Promise(function(resolve, reject) {
          axios.get('/imageLocation/images')
            .then((response) => {
              console.log(response.data);
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
      this.image_Data[items[0].folder_name].forEach(function(filename){
        self.files.push({
          'image': filename,
          'mask': filename
        })
      })
    },

    generateZip () {
      console.log(ZipUtil.getBinaryContent)
      var zip = new Zip()
      let downloadedFiles = []
      for (let folder in this.mask_Data) {
        let folderZip = zip.folder(folder)
        this.mask_Data[folder].forEach(function (filename) {
          var url = `/masks/${folder}/${filename}`
          downloadedFiles.push(
            new Promise(function(resolve, reject) {
              ZipUtil.getBinaryContent(url, function (err, data) {
                if (err) {
                  reject(err)
                }
                folderZip.file(filename, data, {binary: true})
                resolve()
              })
            })
          )
        })
      }
      console.log(downloadedFiles)
      Promise.all(downloadedFiles).then(() => {
        zip.generateAsync({type: 'blob'}).then(function (blob) {
          FileSaver.saveAs(
            blob, 
            new Date(Date.now()).toLocaleDateString().split('/').join('-')
          )
        })
      })
    }
  },
  mounted () {
    this.getFiles()
    this.getFolders()
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
.card-main {
  margin-top: 20px;
  background-color: #343a40;
  border: none;
}
</style>