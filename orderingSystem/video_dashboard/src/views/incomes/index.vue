<template>
  <div class="app-container">
    <!-- <div class="filter-container" style="display: flex;justify-content: space-between;">
      <el-button class="filter-item" type="primary" icon="el-icon-plus" @click="createvideo">
        新建产品
      </el-button>
    </div> -->
    <el-table v-loading="listLoading" :data="orderslist" border fit highlight-current-row
      style="width: 100%;margin-top: 20px;">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="订单类型">
        <template slot-scope="{row}">
          <span>{{ row.ordertype }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="订单号">
        <template slot-scope="{row}">
          <span>{{ row.OrderID }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="收益类型">
        <template slot-scope="{row}">
          <span>{{ row.incometype }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="金额">
        <template slot-scope="{row}">
          <span>{{ row.incomenumber }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="订单时间">
        <template slot-scope="{row}">
          <span>{{ row.createT }}</span>
        </template>
      </el-table-column>

    </el-table>
    <pagination
      :total="totalNum"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.page_size"
      @pagination="getList"
    />

  </div>
</template>

<script>
import { getincomes } from '@/api/table'
import Pagination from '@/components/Pagination'
import store from '@/store'
export default {
  components: { Pagination },
  name: 'incomesinfo',

  data() {
    return {
      listLoading: true,
      data: null,
      listQuery: {
        user_id: "",
        page:1,
        page_size:10
      },
      totalNum: 0,
      dialogVisible: false,
      orderslist: []
    }
  },
  created() {
    this.getList()
  },
  methods: {
    async getList() {
      this.listLoading = true
      this.listQuery.user_id = String(store.getters.id)
      const { data } = await getincomes(this.listQuery)
      this.orderslist = data.incomelist.map(v => {
        this.$set(v, 'edit', false)
        return v
      })
      this.totalNum = data.count
      this.listLoading = false
    },

    createvideo() {
      this.dialogVisible = true
    },
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
