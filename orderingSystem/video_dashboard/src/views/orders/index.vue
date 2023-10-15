<template>
  <div class="app-container">
    <div class="filter-container" style="display: flex;justify-content: space-between;">
      <div>
        <el-input v-model="listQuery.order_id" placeholder="订单号" style="width: 200px;" type="number" class="filter-item"
          @keyup.enter.native="getList" />
        <el-input v-model="listQuery.t_user_id" placeholder="用户id" style="width: 200px;" type="number" class="filter-item"
          @keyup.enter.native="getList" />
        <el-col :span="5">
          <el-date-picker v-model="listQuery.start_date" type="date" placeholder="开始日期" style="width: 100%;" />
        </el-col>
        <el-col :span="5">
          <el-date-picker v-model="listQuery.end_date" type="date" placeholder="结束日期" style="width: 100%;" />
        </el-col>
        <el-button class="filter-item" type="primary" icon="el-icon-search" @click="getList">
          查询
        </el-button>
      </div>
    </div>
    <el-table v-loading="listLoading" :data="orderslist" border fit highlight-current-row
      style="width: 100%;margin-top: 20px;">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.Id }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="用户ID">
        <template slot-scope="{row}">
          <span>{{ row.UserId }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="用户名">
        <template slot-scope="{row}">
          <span>{{ row.UserName }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="订单类型">
        <template slot-scope="{row}">
          <span>{{ row.OrderType }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="订单状态">
        <template slot-scope="{row}">
          <!-- <span>{{ row.OrderStatus }}</span> -->
          <el-tag :type="row.OrderStatus | statusFilter">
            {{ row.OrderStatus}}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="购买产品">
        <template slot-scope="{row}">
          <span>{{ row.Product }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="价格">
        <template slot-scope="{row}">
          <span>{{ row.Price }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="支付方式">
        <template slot-scope="{row}">
          <span>{{ row.PayType }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="下单时间">
        <template slot-scope="{row}">
          <span>{{ row.Create_time}}</span>
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
import { getorders } from '@/api/table'
import Pagination from '@/components/Pagination'
import store from '@/store'
export default {
  components: { Pagination },
  name: 'ordersinfo',
  filters: {
    statusFilter(status) {
      const statusMap = {
        '已支付': 'success',
        '已下单': 'danger'
      }
      return statusMap[status]
    },
  },
  data() {
    return {
      listLoading: true,
      data: null,
      listQuery: {
        user_id: "",
        t_user_id: "",
        start_date: "",
        end_date: "",
        order_id: "",
        page:1,
        page_size:10
      },
      totalNum: 0,
      dialogVisible: false,
      t_user_id: "",
      order_id: "",
      orderslist: []
    }
  },
  created() {
    this.getList()
  },
  methods: {
    formatDate(date) {
      let year = date.getFullYear();
      let month = String(date.getMonth() + 1).padStart(2, '0');
      let day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    },
    async getList() {
      this.listLoading = true
      this.listQuery.user_id = String(store.getters.id)
      if (this.listQuery.end_date != "" && typeof this.listQuery.end_date === "object" && this.listQuery.end_date !== null) {
        this.listQuery.end_date = this.formatDate(this.listQuery.end_date)
      }
      if (this.listQuery.start_date != "" && typeof this.listQuery.start_date === "object" && this.listQuery.start_date !== null) {
        this.listQuery.start_date = this.formatDate(this.listQuery.start_date)
      }
      const { data } = await getorders(this.listQuery)
      console.log(data)
      this.orderslist = data.orderlist
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
