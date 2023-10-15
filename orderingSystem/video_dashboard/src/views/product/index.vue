<template>
  <div class="app-container">
    <div class="filter-container" style="display: flex;justify-content: space-between;">
      <el-button class="filter-item" type="primary" icon="el-icon-plus" @click="createvideo">
        新建产品
      </el-button>
    </div>
    <el-table v-loading="listLoading" :data="productList" border fit highlight-current-row
      style="width: 100%;margin-top: 20px;">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="产品名">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.name" size="small"  />
          </template>
          <span v-else>
            <span>{{ row.name }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column width="200px" align="center" label="产品价格">
        <template slot-scope="{row}">
          <template v-if="row.edit">
            <el-input v-model="row.price" size="small" type="number" />
          </template>
          <span v-else>
            <span>{{ row.price }}</span>
          </span>
        </template>
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
      <el-form label-width="80px" label-position="left">
        <el-form-item label="产品名">
          <el-input placeholder="产品名称" v-model="productName" />
        </el-form-item>
        <el-form-item label="产品价格">
          <el-input placeholder="产品价格" v-model="productPrice" type="number" />
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
import { getProducts, createProducts, deleteProducts } from '@/api/table'
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
      listLoading: true,
      data: null,
      listQuery: {
        page: 1,
        page_size: 20
      },
      totalNum: 0,
      dialogVisible: false,
      productName: "",
      productPrice: "",
      productList: []
    }
  },
  created() {
    this.getList()
  },
  methods: {

    async getList() {
      this.listLoading = true
      const { data } = await getProducts()
      this.productList = data.map(v => {
        this.$set(v, 'edit', false)
        return v
      })
      this.totalNum = data.count
      this.listLoading = false
      console.log(this.productList)
    },
    handleDelete(row) {
      this.$confirm('确认删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          await deleteProducts({ 'pid': row.id })
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
      let data = {}
      if (this.productName  && this.productPrice) {
        try {
          data.name = this.productName
          data.price = parseFloat(this.productPrice)
          await createProducts(data)
          this.getList()
          this.dialogVisible = false
        } catch (error) {
          this.$message({
            type: 'error',
            message: '创建失败'
          })
        }

      } else {
        this.$message({
          type: 'error',
          message: '产品信息有误'
        })
      }
    },
    async confirmEdit(row) {
      try {
        console.log(row)
        let data = {}
        if (row.name  &&  row.price) {
          try {
            data.productid = row.id
            data.name = row.name
            data.price = parseFloat(row.price)
            await createProducts(data)
            this.getList()
            this.dialogVisible = false
          } catch (error) {
            this.$message({
              type: 'error',
              message: '创建失败'
            })
          }

        } else {
          this.$message({
            type: 'error',
            message: '产品信息有误'
          })
        }

        row.edit = !row.edit

      } catch (err) {
        this.$message({
          type: 'error',
          message: '编辑失败'
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
