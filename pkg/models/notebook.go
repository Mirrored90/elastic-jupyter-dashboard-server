package models

import (
	"context"
	"errors"
	"log"

	kubeflowtkestackiov1alpha1 "github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/elastic-jupyter-dashboard-server/pkg/driver"
)

var cl client.Client = driver.MyK8S

type Notebook struct {
	Name             string   `json:"name" binding:"required"`
	Namespace        string   `json:"namespace" binding:"required"`
	PodName          string   `json:"podName"`
	NodeName         string   `json:"node"`
	Label            []string `json:"label"`
	CreatedOn        string   `json:"createdOn"`
	Status           string   `json:"status"`
	GatewayName      string   `json:"gatewayName"`
	GatewayNamespace string   `json:"gatewayNamespace"`
	ContainerName    string   `json:"containerName"`
}

func (model *Notebook) GetNotebooks() ([]Notebook, error) {
	pods := &corev1.PodList{}
	keys := []string{"notebook"}
	if err := cl.List(context.Background(), pods, client.HasLabels(keys)); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	notebookCRs := &kubeflowtkestackiov1alpha1.JupyterNotebookList{}
	if err := cl.List(context.Background(), notebookCRs); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	notebooks := []Notebook{}
	for _, value := range pods.Items {
		notebookName := value.Labels["notebook"]
		notebookNamespace := value.Namespace
		gatewayName, gatewayNamespace := getGatewayInfo(notebookName, notebookNamespace, notebookCRs)

		notebooks = append(notebooks, Notebook{
			Name:             notebookName,
			PodName:          value.Name,
			Namespace:        notebookNamespace,
			NodeName:         value.Spec.NodeName,
			Label:            getLabels(value.Labels),
			CreatedOn:        value.CreationTimestamp.String(),
			Status:           string(value.Status.Phase),
			GatewayName:      gatewayName,
			GatewayNamespace: gatewayNamespace,
		})
	}

	return notebooks, nil
}

func (model *Notebook) DeleteNotebook(name string, namespace string) error {
	notebook := &kubeflowtkestackiov1alpha1.JupyterNotebook{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := cl.Delete(context.Background(), notebook); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (model *Notebook) CreateNotebook() error {
	var notebook *kubeflowtkestackiov1alpha1.JupyterNotebook
	var err error

	notebook, err = model.formatNotebook()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err = cl.Create(context.Background(), notebook); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func getGatewayInfo(name string, namespace string, notebookCRs *kubeflowtkestackiov1alpha1.JupyterNotebookList) (string, string) {
	for _, CR := range notebookCRs.Items {
		if CR.Name == name && CR.Namespace == namespace {
			if CR.Spec.Gateway == nil {
				return "", ""
			}
			return CR.Spec.Gateway.Name, CR.Spec.Gateway.Namespace
		}
	}
	return "", ""
}

func getLabels(labelMap map[string]string) []string {
	labels := []string{}
	for k, v := range labelMap {
		label := k + "=" + v
		labels = append(labels, label)
	}

	return labels
}

func (model *Notebook) formatNotebook() (*kubeflowtkestackiov1alpha1.JupyterNotebook, error) {
	if model.GatewayName == "" && model.ContainerName == "" {
		return &kubeflowtkestackiov1alpha1.JupyterNotebook{}, errors.New("both gateway and container are empty")
	}

	gatewayReference := corev1.ObjectReference{}
	// template := corev1.PodTemplateSpec{}

	if model.GatewayName != "" {
		gateway := &kubeflowtkestackiov1alpha1.JupyterGateway{}

		if err := cl.Get(context.Background(), client.ObjectKey{
			Namespace: model.GatewayNamespace,
			Name:      model.GatewayName,
		}, gateway); err != nil {
			return &kubeflowtkestackiov1alpha1.JupyterNotebook{}, errors.New(err.Error())
		} else {
			gatewayReference = corev1.ObjectReference{
				Name:      model.GatewayName,
				Namespace: model.GatewayNamespace,
			}
		}
	}

	// TODO: fill template if container not empty

	notebook := &kubeflowtkestackiov1alpha1.JupyterNotebook{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: model.Namespace,
		},
		Spec: kubeflowtkestackiov1alpha1.JupyterNotebookSpec{
			Gateway: &gatewayReference,
			// Template: &template,
		},
	}
	return notebook, nil
}
