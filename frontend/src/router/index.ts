import { createRouter, createWebHistory } from "vue-router";

const routes = [
	{
		path: "/login",
		component: () => import("@views/Login.vue"),
		meta: { hideNav: true },
	},
	{ path: "/about", component: () => import("@components/HelloWorld.vue") },
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

export default router;
