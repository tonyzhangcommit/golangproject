<template>
  <div class="app-container">
    <div class="filter-container" style="display: flex;justify-content: space-between;">
      <el-button class="filter-item" type="primary" icon="el-icon-plus" @click="createvideo">
        新建分类
      </el-button>
    </div>
    <el-table v-loading="listLoading" :data="categorylist" border fit highlight-current-row
      style="width: 100%;margin-top: 20px;">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="分类名">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.name" size="small" />
          </template>
          <span v-else>
            <span>{{ row.name }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column width="200px" align="center" label="父级分类">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.categories" size="small" />
          </template>
          <span v-else>
            <span>{{ row.parent || "空" }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="操作" width="230">
        <template slot-scope="{row}">
          <el-button size="mini" type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog :visible.sync="dialogVisible">
      <el-form label-width="80px" label-position="left">
        <el-form-item label="一级分类">
          <div style="display:flex;">
            <el-select v-model="firstlevel" placeholder="一级分类" clearable style="width: 190px" class="filter-item">
              <el-option v-for="item in categoryFL" :key="item.name" :label="item.name" :value="item.name" />
            </el-select>
            <el-input placeholder="自定义" v-model="firstlevel" />
          </div>
        </el-form-item>
        <el-form-item label="二级分类">
          <div style="display:flex;">
            <el-select v-model="secondlevel" placeholder="二级分类" clearable style="width: 190px" class="filter-item">
              <el-option v-for="item in categorySL" :key="item.name" :label="item.name" :value="item.name" />
            </el-select>
            <el-input placeholder="自定义" v-model="secondlevel" />
          </div>
        </el-form-item>
        <el-form-item label="三级分类">
          <div style="display:flex;">
            <el-select v-model="categoryinfoT" placeholder="三级分类" clearable style="width: 190px" class="filter-item">
              <el-option v-for="item in categoryTL" :key="item.name" :label="item.name" :value="item.name" />
            </el-select>
            <el-input placeholder="自定义" v-model="categoryinfoT" />
          </div>
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createVideo">确认</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getcategory, createCategories, deleteCategories } from '@/api/table'
import Pagination from '@/components/Pagination'
export default {
  components: { Pagination },
  name: 'InlineEditTable',
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
      categorylist: [],
      listLoading: true,
      data: null,
      listQuery: {
        page: 1,
        page_size: 20
      },
      totalNum: 0,
      dialogVisible: false,
      categoryFL: [],
      firstlevel: "",
      categorySL: [],
      secondlevel: "",
      categoryTL: [],
      categoryinfoT: "",
    }
  },
  watch: {
    firstlevel(newCategory) {
      try {
        const foundCategory = this.findCategoryByName(this.categoryFL, newCategory);
        this.secondlevel = ""
        this.categorySL = []
        this.categorySL = foundCategory.Categorylist
      } catch {
        console.log("no category")
      }
    },
    secondlevel(newCategory) {
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
  created() {
    this.getList()
    this.InitCategory()
  },
  methods: {
    flattenCategories(categories) {
      const result = [];
      function traverse(category, parentCategory) {
        result.push({
          id: category.id,
          name: category.name,
          intro: category.intro,
          parent: parentCategory ? parentCategory.name : null
        });

        if (category.Categorylist && category.Categorylist.length > 0) {
          category.Categorylist.forEach(subCategory => {
            traverse(subCategory, category);
          });
        }
      }
      categories.forEach(category => {
        traverse(category, null);
      });

      return result;
    },
    async getList() {
      this.listLoading = true
      const { data } = await getcategory()
      let resultList = this.flattenCategories(data)
      this.categorylist = resultList.map(v => {
        this.$set(v, 'edit', false)
        return v
      })
      this.totalNum = data.count
      this.listLoading = false
    },
    handleDelete(row) {
      this.$confirm('确认删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          await deleteCategories({ 'id': row.id })
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
    },
    isValidCategory(firstL, secondL, thirdL) {
      if (thirdL) {
        return firstL && secondL;
      } else if (secondL) {
        return firstL;
      } else {
        return true;
      }
    },
    async createVideo() {
      let data = {}
      if (this.isValidCategory(this.firstlevel, this.secondlevel, this.categoryinfoT)) {
        try {
          console.log(this.firstlevel, this.secondlevel, this.categoryinfoT)
          data.firstlevel = this.firstlevel
          data.secondlevel = this.secondlevel
          data.thirdlevel = this.categoryinfoT
          console.log(data)
          await createCategories(data)
          this.dialogVisible = false
          this.getList()
        } catch (error) {
          this.$message({
            type: 'error',
            message: '创建失败'
          })
        }

      } else {
        this.$message({
          type: 'error',
          message: '请填写正确分类信息'
        })
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
