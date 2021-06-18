package routers

import (
	"github.com/gin-gonic/gin"
	"mana/src/controllers/auth"
	"mana/src/controllers/kubernetes"
	"mana/src/controllers/menus"
	"mana/src/controllers/navigation"
	"mana/src/controllers/resource"
	"mana/src/controllers/user"
	"mana/src/filters/util"
	"net/http"
)

// Index 路由总配置
func Index(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"code": "200", "message": "successfully"}) })
	r.HEAD("/", func(c *gin.Context) { c.JSON(http.StatusOK, "successfully") })

	// 用户登录,注册，
	r.POST("/api/v1/common/user/login", user.Login)
	r.POST("/api/v1/common/user/register", user.Login)

	// 验证请求token中间件
	r.Use(util.AuthToken())
	// 路由
	v1Group := r.Group("/api/v1/common")
	{

		// 获取左侧菜单
		v1Group.GET("/menus/list", menus.GetMenus)

		// 系统管理
		v1Group.GET("/system/role", auth.GetRoleList)
		v1Group.POST("/system/role", auth.AddRole)
		v1Group.DELETE("/system/role", auth.DeleteRole)
		v1Group.GET("/system/menus/list", menus.GetMenusAll)
		v1Group.PATCH("/system/menus/role/binding", auth.UpdateRolePermission)

		// 获取用户详情
		userGroup := v1Group.Group("/user")
		{
			userGroup.GET("/userinfo/:uid", user.FindByUserinfo)
		}

		// 获取导航链接列表，添加链接，编辑，删除，获取单条导航链接记录
		navGroup := v1Group.Group("/navigation")
		{
			navGroup.GET("/links", navigation.GetResourceLinks)
			navGroup.POST("/links", navigation.AddResourceLink)
			navGroup.GET("/links/:id", navigation.GetResourceLinks)
			navGroup.PATCH("/links/:id", navigation.UpdateResourceLink)
			navGroup.DELETE("/links/:id", navigation.DeleteResourceLink)
		}

		// 主机监控
		v1Group.GET("/resource/monitor", resource.GetHostMonitorInfo)

		// k8s 相关
		k8sGroup := v1Group.Group("/kubernetes")
		{
			// 添加k8s配置文件
			k8sGroup.POST("/cluster", kubernetes.InstKubeConfig)
			// 获取集群配置列表
			k8sGroup.GET("/cluster", kubernetes.GetKubeConfig)
			// 删除集群配置
			k8sGroup.DELETE("/cluster/:cid", kubernetes.DelKubeConfig)
			// 获取k8s名称空间/common/kubernetes/cluster/:cid/work/namespaces
			k8sGroup.GET("/cluster/:cid/namespaces", kubernetes.GetNamespace)
			// 获取k8s控制器资源
			k8sGroup.GET("/cluster/:cid/control/:namespaces/:control", kubernetes.GetKubernetesControl)
			// 获取k8s工作负载详细信息
			// /apis/apps/v1/namespaces/kube-system/daemonsets/kube-flannel-ds-amd64
			k8sGroup.GET("/cluster/:cid/pods/:namespaces/:control/:podsName", kubernetes.GetKubernetesPods)
		}

	}
}
