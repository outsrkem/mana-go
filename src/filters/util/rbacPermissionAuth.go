package util

func RbacPermissionAuth(p string, perList *[]string) bool {
    // perList = [menu:link:add menu:link:update]
    // p = "menu:link:update"
    for _, v := range *perList {
        if p == v {
            _log.Info("权限校验通过",p ," ==> ", *perList)
            return true
        }
    }
    _log.Info("权限校验失败",p ," ==> ", *perList)
    return false
}
