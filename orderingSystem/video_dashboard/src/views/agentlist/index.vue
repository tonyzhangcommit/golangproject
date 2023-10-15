<template>
  <div class="app-container">
    <div
      class="filter-container"
      style="display: flex; justify-content: space-between"
    >
      <div>
        <el-input
          v-model="listQuery.name"
          placeholder="用户名"
          style="width: 200px"
          class="filter-item"
          @keyup.enter.native="getList"
        />
        <el-input
          v-model="listQuery.veriCode"
          placeholder="识别码"
          style="width: 200px"
          class="filter-item"
          @keyup.enter.native="getList"
        />
        <el-input
          v-model="listQuery.telNum"
          placeholder="手机号"
          style="width: 200px"
          class="filter-item"
          @keyup.enter.native="getList"
        />
        <el-button
          class="filter-item"
          type="primary"
          icon="el-icon-search"
          @click="getList"
        >
          查询
        </el-button>
      </div>
      <el-button
        class="filter-item"
        type="primary"
        icon="el-icon-plus"
        @click="createuser"
      >
        新建用户
      </el-button>
    </div>
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%; margin-top: 20px"
    >
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{ row }">
          <span>{{ row.ID }}</span>
        </template>
      </el-table-column>
      <el-table-column width="180px" align="center" label="用户名">
        <template slot-scope="{ row }">
          <span>{{ row.Name }}</span>
        </template>
      </el-table-column>
      <el-table-column width="120px" align="center" label="识别码">
        <template slot-scope="{ row }">
          <span>{{ row.IdentificationCode }}</span>
        </template>
      </el-table-column>
      <el-table-column width="200px" align="center" label="手机号">
        <template slot-scope="{ row }">
          <span>{{ row.Telnumber }}</span>
        </template>
      </el-table-column>
      <el-table-column width="120px" align="center" label="上级代理">
        <template slot-scope="{ row }">
          <span>{{ row.ManagerID }}</span>
        </template>
      </el-table-column>
      <el-table-column width="120px" align="center" label="子代理数量">
        <template slot-scope="{ row }">
          <span>{{ row.CountSelfUser }}</span>
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="用户状态" min-width="110">
        <template slot-scope="{ row }">
          <template v-if="row.edit">
            <el-input
              v-model="row.status"
              placeholder="0为封禁1为正常"
              class="edit-input"
              size="small"
              type="number"
            />
            <el-button
              class="cancel-btn"
              size="small"
              icon="el-icon-refresh"
              type="warning"
              @click="cancelEdit(row)"
            >
              取消
            </el-button>
          </template>
          <span v-else>
            <el-tag :type="row.Status | statusFilter">
              {{ row.Status | showStatusFilter }}
            </el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column min-width="120" align="center" label="创建时间">
        <template slot-scope="{ row }">
          <span>{{ row.CreateTime }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="操作" width="120">
        <template slot-scope="{ row }">
          <el-button
            v-if="row.edit"
            type="success"
            size="small"
            icon="el-icon-circle-check-outline"
            @click="confirmEdit(row)"
          >
            确认
          </el-button>
          <el-button
            v-else
            type="primary"
            size="small"
            icon="el-icon-edit"
            @click="row.edit = !row.edit"
          >
            编辑
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog :visible.sync="dialogVisible">
      <el-form label-width="80px" label-position="left">
        <el-form-item label="用户名">
          <el-input placeholder="用户名" v-model="userItem.name" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input placeholder="密码" v-model="userItem.passWord" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input placeholder="手机号" v-model="userItem.mobile" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select
            v-model="userItem.role"
            placeholder="角色"
            clearable
            style="width: 190px"
            class="filter-item"
          >
            <el-option
              v-for="item in roleslist"
              :key="item.id"
              :label="item.name"
              :value="item.name"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div style="text-align: right">
        <el-button type="danger" @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createUser">确认</el-button>
      </div>
    </el-dialog>
    <pagination
      :total="totalNum"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.page_size"
      @pagination="getList"
    />
  </div>
</template>

<script>
import {
  getuserlist,
  getroles,
  createuserapi,
  edituserstatus,
} from "@/api/table";
import Pagination from "@/components/Pagination";
import store from "@/store";
export default {
  components: { Pagination },
  name: "InlineEditTable",
  filters: {
    statusFilter(status) {
      const statusMap = {
        true: "success",
        false: "danger",
      };
      return statusMap[status];
    },
    showStatusFilter(status) {
      const statusMap = {
        true: "正常",
        false: "封禁",
      };
      return statusMap[status];
    },
  },
  data() {
    return {
      list: null,
      listLoading: true,
      listQuery: {
        id: "",
        name: "",
        veriCode: "",
        telNum: "",
        page: 1,
        page_size: 10,
      },
      totalNum: 0,
      dialogVisible: false,
      userItem: {
        name: "",
        password: "",
        mobile: "",
        role: "",
        managerid: 0,
      },
      getrolesquery: {
        id: "",
      },
      roleslist: [],
      selectRoles: "",
      tempStatus: 0,
      edituser:{
        id:0,
        status:0,
        targetid:0,
      },
    };
  },
  created() {
    this.getList();
    this.getRoles();
  },
  methods: {
    async getList() {
      this.listLoading = true;
      this.listQuery.id = String(store.getters.id);
      const { data } = await getuserlist(this.listQuery);
      if (data.userlist) {
        this.list = data.userlist.map((v) => {
          this.$set(v, "edit", false);
          v.status = v.status;
          return v;
        });
      } else {
        this.list = [];
      }
      this.totalNum = data.count;
      this.listLoading = false;
    },
    cancelEdit(row) {
      row.title = row.originalTitle;
      row.edit = false;
      this.$message({
        message: "取消操作",
        type: "warning",
      });
    },
    async confirmEdit(row) {
      this.edituser.id = store.getters.id;
      this.edituser.status = parseInt(row.status);
      this.edituser.targetid = row.ID;
      try {
        const { data } = await edituserstatus(this.edituser);
        this.$message({
          message: "用户状态已改变",
          type: "success",
        });
        this.getList();
      } catch (error) {
        this.$message({
          type: "error",
          message: error,
        });
      }

      row.edit = false;
      row.originalTitle = row.title;
    },
    handleFilter() {},
    createuser() {
      this.dialogVisible = true;
    },
    async createUser() {
      this.userItem.managerid = store.getters.id;
      try {
        const { data } = await createuserapi(this.userItem);
        this.$message({
          type: "success",
          message: "创建成功",
        });
        this.dialogVisible = false;
        this.getList();
      } catch (error) {
        this.$message({
          type: "error",
          message: data.msg,
        });
      }
    },
    async getRoles() {
      this.getrolesquery.id = String(store.getters.id);
      const { data } = await getroles(this.listQuery);
      this.roleslist = data;
    },
  },
};
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
