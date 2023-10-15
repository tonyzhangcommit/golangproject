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
      component: () => import('@/views/agentlist/index'),
      meta: { title: '代理列表', icon: 'dashboard' }
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
  // {
  //   path: '/agentlist',
  //   component: Layout,
  //   children: [
  //     {
  //       path: 'agentlist',
  //       name: 'agentlist',
  //       component: () => import('@/views/agentlist/index'),
  //       meta: { title: '代理列表', icon: 'form', roles: ['superadmin'] }
  //     }
  //   ]
  // },

  {
    path: '/videomanage',
    component: Layout,
    redirect: '/videomanage/videolist',
    name: 'videomanage',
    meta: { title: '产品管理', icon: 'el-icon-s-help', roles: ['superadmin'] },
    children: [
      {
        path: 'videolist',
        name: 'videolist',
        component: () => import('@/views/videolist/index'),
        meta: { title: '视频列表', icon: 'table', roles: ['superadmin'] }
      },
      {
        path: 'category',
        name: 'category',
        component: () => import('@/views/category/index'),
        meta: { title: '分类管理', icon: 'tree', roles: ['superadmin'] }
      },
      {
        path: 'product',
        name: 'product',
        component: () => import('@/views/product/index'),
        meta: { title: '产品列表', icon: 'tree', roles: ['superadmin'] }
      }
    ]
  },

  // {
  //   path: '/proxylist',
  //   component: Layout,
  //   children: [
  //     {
  //       path: 'proxylist',
  //       name: 'proxylist',
  //       component: () => import('@/views/agentlist/index'),
  //       meta: { title: '代理列表', icon: 'form', roles: ['manager'] }
  //     }
  //   ]
  // },

  {
    path: '/realtimeorder',
    component: Layout,
    children: [
      {
        path: 'realtimeorder',
        name: 'realtimeorder',
        component: () => import('@/views/orders/index'),
        meta: { title: '订单详情', icon: 'form', roles: ['superadmin','manager'] }
      }
    ]
  },
  {
    path: '/readincome',
    component: Layout,
    children: [
      {
        path: 'readincome',
        name: 'readincome',
        component: () => import('@/views/incomes/index'),
        meta: { title: '收益详情', icon: 'form', roles: ['superadmin'] }
      }
    ]
  },

  // {
  //   path: '/shoperinfo',
  //   component: Layout,
  //   children: [
  //     {
  //       path: 'shoperinfo',
  //       name: 'shoperinfo',
  //       component: () => import('@/views/form/index'),
  //       meta: { title: '个人信息', icon: 'form', roles: ['superadmin'] }
  //     }
  //   ]
  // },



  {
    path: '/incomeinfo',
    component: Layout,
    children: [
      {
        path: 'incomeinfo',
        name: 'incomeinfo',
        component: () => import('@/views/incomes/index'),
        meta: { title: '收益详情', icon: 'form', roles: ['manager'] }
      }
    ]
  },

  // {
  //   path: '/managerinfo',
  //   component: Layout,
  //   children: [
  //     {
  //       path: 'managerinfo',
  //       name: 'managerinfo',
  //       component: () => import('@/views/form/index'),
  //       meta: { title: '个人资料', icon: 'form', roles: ['manager'] }
  //     }
  //   ]
  // },

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
