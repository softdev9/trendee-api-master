package appinvite

import (
	"fmt"
	"net/http"
)

// GetAppInviteHTML Builds an app invite page
func GetAppInviteHTML(rw http.ResponseWriter, req *http.Request) {
	var appInviteContent = `<html>
<head>
    <meta property="al:ios:url" content="trendee://" />
    <meta property="al:ios:app_store_id" content="979043211" />
    <meta property="al:ios:app_name" content="Trendee App" />
    <meta property="al:android:url" content="trendee://" />
    <meta property="al:android:app_name" content="Trendee" />
    <meta property="al:android:package" content="com.trendeelab.trendee" />
</head>
<body>
Link to the app
</body>
</html>`
	fmt.Fprint(rw, appInviteContent)

}
