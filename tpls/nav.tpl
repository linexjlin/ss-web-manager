{{define "nav"}}
            <nav class="">
                <div style="overflow:hidden; z-index: inherit;" class="navbar king-horizontal-nav1 king-horizontal-primary f14">
                    <div class="navbar-container">
                        <div class="navbar-header pull-left">
                            <a class="navbar-brand" href="javascript:;">
                                <img src="https://o.qcloud.com/static_api/v3/bk/images/logo.png" class="logo"> </a>
                        </div>
                        <ul class="nav navbar-nav pull-left m0">
                            <li class="inactive"><a href="/login">登录</a></li>
                            <li class="inactive"><a href="/signup">注册</a></li>
                            <li class="inactive"><a href="/us">联系我们</a></li>
                        </ul>
                        <div class="navbar-header pull-right">
                            <ul class="nav">
                                <li class="user-info">
                                    <a href="javascript:;">
                                        <img class="img-rounded" src="https://o.qcloud.com/static_api/v3/components/horizontal_nav1/images/avatar.png">
					<span>{{.Name}}</span>
                                    </a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </nav>
{{end}}
