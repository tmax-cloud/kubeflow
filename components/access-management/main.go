/*
 * Kubeflow Auth
 *
 * Access Management API.
 *
 * API version: 0.1.0
 * Contact: kubeflow-engineering@google.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"flag"
	"net/http"
	"k8s.io/klog"
//	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/kubeflow/kubeflow/components/access-management/kfam"
	profile "github.com/kubeflow/kubeflow/components/access-management/pkg/apis/kubeflow/v1beta1"

	istioSecurityClient "istio.io/client-go/pkg/apis/security/v1beta1"
)

// kfam API assume coming request will contain user id in request header.
// set this parameter to specify header key containing user id.
const USERIDHEADER = "userid-header"

// set this parameter to specify header value prefix (if any) before user id.
const USERIDPREFIX = "userid-prefix"

// set cluster admin user id here.
const CLUSTERADMIN = "cluster-admin"

func main() {		

	klog.V(3).Infof("Server started")
	var LogLevel string
	var userIdHeader string
	var userIdPrefix string
	var clusterAdmin string
	flag.StringVar(&LogLevel, "log-level", "INFO", "Log Level; INFO, WARN, ERROR, FATAL")
	flag.StringVar(&userIdHeader, USERIDHEADER, "x-goog-authenticated-user-email", "Key of request header containing user id")
	flag.StringVar(&userIdPrefix, USERIDPREFIX, "accounts.google.com:", "Request header user id common prefix")
	flag.StringVar(&clusterAdmin, CLUSTERADMIN, "", "cluster admin")
	flag.Parse()

	if LogLevel == "INFO" || LogLevel == "info" {
        LogLevel = "3"
    } else if LogLevel == "WARN" || LogLevel == "warn" {
        LogLevel = "2"
    } else if LogLevel == "ERROR" || LogLevel == "error" {
        LogLevel = "1"
    } else if LogLevel == "FATAL" || LogLevel == "fatal" {
        LogLevel = "0"
    } else {
        klog.Infoln("Unknown log-level paramater. ")
        LogLevel = "3"
    }
//    klog.InitFlags(nil)
    flag.Set("v", LogLevel)
	profile.AddToScheme(scheme.Scheme)
	istioSecurityClient.AddToScheme(scheme.Scheme)

	profileClient, err := kfam.NewKfamClient(userIdHeader, userIdPrefix, clusterAdmin)
	if err != nil {
		klog.V(1).Info(err)
		panic(err)
	}

	router := kfam.NewRouter(profileClient)

	klog.V(0).Info(http.ListenAndServe(":8081", router))
}
