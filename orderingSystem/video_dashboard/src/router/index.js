import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [{
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboardmanager/index'),
      meta: { title: '首页', icon: 'dashboard' }
    }]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  }
]

// manager superadmin

export const asyncRoutes = [
  {
    path: '/shopmanage',
    component: Layout,
    redirect: '/shopmanage/shoporder',
    name: 'shopmanage',
    meta: { title: '店铺管理', icon: 'el-icon-s-help', roles: ['superadmin'] },
    children: [
      {
        path: 'shoporder',
        name: 'shoporder',
        component: () => import('@/views/table/index'),
        meta: { title: '订单管理', icon: 'table', roles: ['superadmin'] }
      },
      {
        path: 'menu',
        name: 'menu',
        component: () => import('@/views/tree/index'),
        meta: { title: '菜单管理', icon: 'tree', roles: ['superadmin'] }
      },
      {
        path: 'table',
        name: 'table',
        component: () => import('@/views/tree/index'),
        meta: { title: '桌号管理', icon: 'tree', roles: ['superadmin'] }
      },
      {
        path: 'shopinfo',
        name: 'shopinfo',
        component: () => import('@/views/tree/index'),
        meta: { title: '店铺信息', icon: 'tree', roles: ['superadmin'] }
      }
    ]
  },

  {
    path: '/realtimeorder',
    component: Layout,
    children: [
      {
        path: 'realtimeorder',
        name: 'realtimeorder',
        component: () => import('@/views/form/index'),
        meta: { title: '实时订单', icon: 'form', roles: ['superadmin'] }
      }
    ]
  },

  {
    path: '/shoperinfo',
    component: Layout,
    children: [
      {
        path: 'shoperinfo',
        name: 'shoperinfo',
        component: () => import('@/views/form/index'),
        meta: { title: '个人信息', icon: 'form', roles: ['superadmin'] }
      }
    ]
  },

  {
    path: '/proxylist',
    component: Layout,
    children: [
      {
        path: 'proxylist',
        name: 'proxylist',
        component: () => import('@/views/form/index'),
        meta: { title: '代理列表', icon: 'form', roles: ['manager'] }
      }
    ]
  },

  {
    path: '/incomeinfo',
    component: Layout,
    children: [
      {
        path: 'incomeinfo',
        name: 'incomeinfo',
        component: () => import('@/views/form/index'),
        meta: { title: '收益详情', icon: 'form', roles: ['manager'] }
      }
    ]
  },

  {
    path: '/managerinfo',
    component: Layout,
    children: [
      {
        path: 'managerinfo',
        name: 'managerinfo',
        component: () => import('@/views/form/index'),
        meta: { title: '个人资料', icon: 'form', roles: ['manager'] }
      }
    ]
  },

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
