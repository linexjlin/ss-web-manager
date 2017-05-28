{{define "nav"}}
            <nav class="">
                <div style="overflow:hidden; z-index: inherit;" class="navbar king-horizontal-nav1 king-horizontal-primary f14">
                    <div class="navbar-container">
                        <div class="navbar-header pull-left">
                            <a class="navbar-brand" href="javascript:;">
                                <img src="https://o.qcloud.com/static_api/v3/bk/images/logo.png" class="logo"> </a>
                        </div>
                        <ul class="nav navbar-nav pull-left m0">
                            <li class="inactive"><a href="/login"> <i class="fa fa-sign-in"></i> 登录</a></li>
                            <li class="inactive"><a href="/signup"> <i class="fa fa-pencil-square-o"></i>  注册</a></li>
                            <li class="inactive"><a href="/us"> <i class="fa fa-group"></i>  关于</a></li>
                        </ul>
                        <div class="navbar-header pull-right">
                            <ul class="nav">
                                <li class="user-info">
                                    <a href="javascript:;">
					<i class="fa fa-user"></i>
					<span>{{.Name}}</span>
                                    </a>
				    <a href="/logout">
					<i class="fa fa-sign-out"></i>
				    </a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </nav>
{{end}}
