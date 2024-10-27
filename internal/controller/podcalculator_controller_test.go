package controller

import (
	podv1alpha1 "github.com/zmotso/pod-calculator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("PodCalculator Controller", func() {

	BeforeEach(func() {
		By("Creating pods")

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "nginx",
					Image: "nginx",
				}},
			},
		}
		Expect(k8sClient.Create(ctx, pod)).To(Succeed())

		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod-with-label",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "nginx",
					Image: "nginx",
				}},
			},
		}
		Expect(k8sClient.Create(ctx, pod)).To(Succeed())
	})

	AfterEach(func() {
		pod := &corev1.Pod{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name:      "test-pod",
			Namespace: "default",
		}, pod)
		Expect(err).NotTo(HaveOccurred())

		Expect(k8sClient.Delete(ctx, pod)).To(Succeed())

		pod = &corev1.Pod{}
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      "test-pod-with-label",
			Namespace: "default",
		}, pod)
		Expect(err).NotTo(HaveOccurred())

		Expect(k8sClient.Delete(ctx, pod)).To(Succeed())
	})
	It("Should successfully reconcile the resource", func() {
		By("Creating PodCalculator resource")

		podCalculator := &podv1alpha1.PodCalculator{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod-calculator",
				Namespace: "default",
			},
			Spec: podv1alpha1.PodCalculatorSpec{
				Namespace: "default",
				LabelsSelector: map[string]string{
					"app": "test",
				},
				ConfigmapRef: &podv1alpha1.ConfigMapRef{
					Name:      "test-pod-calculator",
					Namespace: "default",
				},
			},
		}

		Expect(k8sClient.Create(ctx, podCalculator)).To(Succeed())
		Eventually(func(g Gomega) {
			createdUser := &podv1alpha1.PodCalculator{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: "test-pod-calculator", Namespace: "default"}, createdUser)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(createdUser.Status.Count).Should(Equal(1))
		}).WithTimeout(time.Second * 5).Should(Succeed())
	})
})
