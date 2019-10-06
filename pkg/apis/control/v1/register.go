package v1

import(
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kinnylee.com/crds-controller-demo/pkg/apis/control"

)

var	SchemeGroupVersion = schema.GroupVersion{
		Group: control.GroupName,
		Version: "v1",
}

func Resource(resource string) schema.GroupResource{
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var(
	SchemeBuilder runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register()
}

func addKnownTypes(schema *runtime.Scheme) error {
	schema.AddKnownTypes(SchemeGroupVersion,
		&Scaling{},
		&ScalingList{},
	)
	metav1.AddToGroupVersion(schema, SchemeGroupVersion)
	return nil
}