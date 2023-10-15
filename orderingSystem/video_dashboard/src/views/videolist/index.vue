<template>
  <div class="app-container">
    <div class="filter-container" style="display: flex;justify-content: space-between;">
      <div>
        <el-input v-model="listQuery.name" placeholder="剧名" style="width: 200px;" class="filter-item"
          @keyup.enter.native="getList" />
        <el-button class="filter-item" type="primary" icon="el-icon-search" @click="getList">
          查询
        </el-button>
      </div>
      <el-button class="filter-item" type="primary" icon="el-icon-plus" @click="createvideo">
        新建视频
      </el-button>
    </div>
    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%;margin-top: 20px;">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="剧名">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.name" size="small" />
          </template>
          <span v-else>
            <span>{{ row.name }}</span>
          </span>
        </template>
      </el-table-column>

      <el-table-column min-width="120px" align="center" label="剧照">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.cover" size="small" />
          </template>
          <span v-else>
            <span>{{ row.cover }}</span>
          </span>
        </template>
        <!-- <template slot-scope="{row}">
        </template> -->
      </el-table-column>
      <el-table-column width="200px" align="center" label="影视简介">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.intro" size="small" />
          </template>
          <span v-else>
            <span>{{ row.intro }}</span>
          </span>
        </template>
        <!-- <template slot-scope="{row}">
          <span>{{ row.intro }}</span>
        </template> -->
      </el-table-column>
      <el-table-column width="200px" align="center" label="分类">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.categories" size="small" />
          </template>
          <span v-else>
            <span>{{ row.categories }}</span>
          </span>
        </template>
        <!-- <template slot-scope="{row}">
          <span>{{ row.categories }}</span>
        </template> -->
      </el-table-column>
      <el-table-column width="120px" align="center" label="剧集">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.episodes" size="small" />
          </template>
          <span v-else>
            <span>{{ row.episodes }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column min-width="120px" align="center" label="剧集链接">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.url" size="small" />
          </template>
          <span v-else>
            <span>{{ row.url }}</span>
          </span>
        </template>
        <!-- <template slot-scope="{row}">
          <span>{{ row.url }}</span>
        </template> -->
      </el-table-column>
      <el-table-column width="120px" align="center" label="剧集简介">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.videoIntro" size="small" />
          </template>
          <span v-else>
            <span>{{ row.videoIntro }}</span>
          </span>
        </template>
        <!-- <template slot-scope="{row}">
          <span>{{ row.videoIntro }}</span>
        </template> -->
      </el-table-column>
      <el-table-column align="center" label="操作" width="230">
        <template slot-scope="{row}">
          <el-button v-if="row.edit" type="success" size="mini" icon="el-icon-circle-check-outline"
            @click="confirmEdit(row)">
            确认
          </el-button>
          <el-button v-else type="primary" size="mini" icon="el-icon-edit" @click="row.edit = !row.edit">
            编辑
          </el-button>
          <el-button size="mini" type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog :visible.sync="dialogVisible">
      <div style="margin-bottom:30px;">
        <el-radio-group v-model="iscreatevideo">
          <el-radio-button label="新建影视" />
          <el-radio-button label="新建剧集" />
        </el-radio-group>
      </div>
      <el-form label-width="80px" label-position="left">
        <el-form-item label="影视名称" v-if="iscreatevideo == '新建影视'">
          <el-input placeholder="影视名" v-model="videoinfo.name" />
        </el-form-item>
        <el-form-item label="影视简介" v-if="iscreatevideo == '新建影视'">
          <el-input :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" placeholder="影视简介" v-model="videoinfo.intro" />
        </el-form-item>
        <el-form-item label="影视封面" v-if="iscreatevideo == '新建影视'">
          <el-input placeholder="影视封面" v-model="videoinfo.cover" />
        </el-form-item>
        <el-form-item label="分类" v-if="iscreatevideo == '新建影视'">
          <el-select v-model="categoryinfoF" placeholder="一级分类" clearable style="width: 190px" class="filter-item">
            <el-option v-for="item in categoryFL" :key="item.id" :label="item.name" :value="item.name" />
          </el-select>
          <el-select v-model="categoryinfoS" placeholder="二级分类" clearable style="width: 190px" class="filter-item">
            <el-option v-for="item in categorySL" :key="item.id" :label="item.name" :value="item.name" />
          </el-select>
          <el-select v-model="categoryinfoT" placeholder="三级分类" clearable style="width: 190px" class="filter-item">
            <el-option v-for="item in categoryTL" :key="item.id" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="影视名称" v-if="iscreatevideo == '新建剧集'">
          <el-select v-model="newVideoName" placeholder="影视名称" clearable style="width: 190px" class="filter-item">
            <el-option v-for="item in videoNameList" :key="item.index" :label="item.name" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="剧集" v-if="iscreatevideo == '新建剧集'">
          <el-input placeholder="剧集" v-model="episodesitem.episodes" type="number" />
        </el-form-item>
        <el-form-item label="剧集地址" v-if="iscreatevideo == '新建剧集'">
          <el-input placeholder="剧集地址" v-model="episodesitem.url" />
        </el-form-item>
        <el-form-item label="剧集简介" v-if="iscreatevideo == '新建剧集'">
          <el-input :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" placeholder="剧集简介"
            v-model="episodesitem.intro" />
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createVideo">确认</el-button>
      </div>
    </el-dialog>
    <pagination :total="totalNum" :page.sync="listQuery.page" :limit.sync="listQuery.page_size" @pagination="getList" />
  </div>
</template>

<script>
import { getvideolist, deletevideo, getcategory, createVideo, createEpisodes } from '@/api/table'
import Pagination from '@/components/Pagination'
export default {
  components: { Pagination },
  name: 'InlineEditVideo',
  filters: {
    statusFilter(status) {
      const statusMap = {
        true: 'success',
        false: 'danger'
      }
      return statusMap[status]
    },
    showStatusFilter(status) {
      const statusMap = {
        true: '正常',
        false: '封禁'
      }
      return statusMap[status]
    },
  },
  data() {
    return {
      list: [],
      videoNameList: [],
      listLoading: true,
      data: null,
      listQuery: {
        page: 1,
        page_size: 10
      },
      categoryQuery: {
        categoryid: 0
      },
      totalNum: 0,
      dialogVisible: false,
      iscreatevideo: "新建影视",
      videoinfo: {
        name: "",
        intro: "",
        cover: ""
      },
      episodesitem: {
        episodes: "",
        url: "",
        intro: "",
      },
      categoryFL: [],
      categoryinfoF: "",
      categorySL: [],
      categoryinfoS: "",
      categoryTL: [],
      categoryinfoT: "",
      newVideoName: ""

    }
  },
  created() {
    this.getList()
    this.InitCategory()
  },
  watch: {
    categoryinfoF(newCategory) {
      try {
        const foundCategory = this.findCategoryByName(this.categoryFL, newCategory);
        this.categoryinfoS = ""
        this.categorySL = []
        this.categorySL = foundCategory.Categorylist
      } catch {
        console.log("no category")
      }
    },
    categoryinfoS(newCategory) {
      try {
        const foundCategory = this.findCategoryByName(this.categorySL, newCategory);
        this.categoryinfoT = ""
        this.categoryTL = []
        this.categoryTL = foundCategory.Categorylist
      } catch {
        console.log("no category")
      }
    }
  },
  methods: {
    async getList() {
      this.listLoading = true
      const { data } = await getvideolist(this.listQuery)
      let videos = data.videos
      let templist = []
      videos.forEach(item => {
        if (item.Videolist.length == 0) {
          templist.push(item)
        } else {
          item.Videolist.forEach(video => {
            templist.push({
              id: item.id,
              name: item.name,
              cover: item.cover,
              intro: item.intro,
              categories: item.Categories.map(category => category.name).join('-'),
              videoID: video.id,
              episodes: video.episodes,
              url: video.url,
              videoIntro: video.intro
            });
          });
        }
      });
      this.list = templist.map(v => {
        if (!this.videoNameList.includes(v.name)) {
          this.videoNameList.push(v.name)
        }
        this.$set(v, 'edit', false)
        return v
      })
      this.totalNum = data.count
      this.listLoading = false
    },
    confirmEdit(row) {
      row.edit = false
      this.$confirm('确认删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          await deletevideo({ 'videoinfo': row })
          this.$message({
            type: 'success',
            message: '视频信息已更改'
          })
        })
        .catch(err => { console.error(err) })
    },
    handleDelete(row) {
      console.log(this.list)
      this.$confirm('确认删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          await deletevideo({ 'videoid': parseInt(row.id), 'videoitemid': parseInt(row.videoID) })
          this.$message({
            type: 'success',
            message: '删除成功！'
          })
          this.getList()
        })
        .catch(err => { console.error(err) })
    },
    createvideo() {
      this.dialogVisible = true
    },
    async createVideo() {
      let categories = []
      if (this.iscreatevideo == '新建影视') {
        if (this.categoryinfoF == "" || this.categoryinfoS == "") {
          this.$message({
            type: 'error',
            message: '请选择一/二级别分类'
          })
        } else {
          try {
            categories.push(this.findCategoryIdByName(this.categoryFL, this.categoryinfoF))
            categories.push(this.findCategoryIdByName(this.categorySL, this.categoryinfoS))
            if (this.categoryinfoT) {
              categories.push(this.findCategoryIdByName(this.categoryTL, this.categoryinfoT))
            }
            this.videoinfo.categories = categories
            await createVideo(this.videoinfo)
          } catch (err) {
            this.$message({
              type: 'error',
              message: '创建失败'
            })
          }
        }
      } else {
        try {
          let videoId = 0
          for (const item of this.list) {
            if (item.name === this.newVideoName) {
              videoId = item.id
              break
            }
          }
          this.episodesitem.episodes = parseInt(this.episodesitem.episodes)
          this.episodesitem.videoid = videoId
          console.log(this.episodesitem)
          await createEpisodes(this.episodesitem)
        } catch (err) {
          this.$message({
            type: 'error',
            message: '创建失败'
          })
        }
      }
      this.dialogVisible = false
      this.getList()
    },
    async InitCategory() {
      const { data } = await getcategory()
      this.categoryFL = data
    },
    findCategoryByName(categoryList, categoryName) {
      for (const category of categoryList) {
        if (category.name === categoryName) {
          return category;
        }
      }
    },
    findCategoryIdByName(categoryList, categoryName) {
      for (const category of categoryList) {
        if (category.name === categoryName) {
          return category.id;
        }
      }
    }
  }
}
</script>

<style scoped>
.edit-input {
  padding-right: 100px;
}

.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
</style>
