// this middleware is for blocking bots from trying to access sensitive files
// I am sick of these bots, in the future if we have time, we can set up a load balancer and attach firewalls but this will do for now
package middlewares

import (
	"net/http"
	"strings"
)

func BlockSuspiciousRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// list of suspicious paths from chatgpt, seems robust enough for now
		suspiciousPaths := []string{
			"phpinfo.php", "admin", "test.php", "debug.php", ".php",
			"_profiler/phpinfo", "_profiler/phpinfo.php", "backend/phpinfo.php",
			"server-info.php", "server_info.php", "config.php", "configuration.php",
			".env", "secrets", "secret", ".git", ".gitignore", ".htaccess", ".htpasswd",
			"wp-config.php", "config.json", "config.yaml", "config.xml",
			"database.json", "database.yaml", "database.xml",
			"settings.php", "settings.json", "settings.xml",
			"admin/", "wp-admin", "wp-login.php", "login.php", "dashboard",
			"cms/", "backend/", "console/", "system/", "cpanel/", "webadmin",
			"administrator/", "user/", "manager/", "root/", "superadmin",
			"debug.php", "debug/", "log", "logs", "error.log", "access.log",
			"trace.log", "console.log", "docker-compose.yml", "composer.json",
			"apikey", "api_key", "token", "jwt", "oauth",
			"backup", "backup.zip", "backup.tar", "backup.sql",
			"old/", "old-site/", "old_version/", "archive/", "bkp/",
			"xmlrpc.php", "shell.php", "cmd.php", "eval.php", "exec.php", "payload.php",
			"wp-content/plugins", "wp-content/uploads", "phpMyAdmin", "mysql-admin",
			"sql/", "phpmyadmin/", "phpmyadmin/setup/",
			"uploads/", "tmp/", "temp/", "cache/", "assets/",
		}

		for _, path := range suspiciousPaths {
			if strings.Contains(r.URL.Path, path) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
