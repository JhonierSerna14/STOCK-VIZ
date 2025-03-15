import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import HomeView from "../views/HomeView.vue";
import StocksView from "../views/StocksView.vue";
import RecommendationsView from "../views/RecommendationsView.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "home",
    component: HomeView,
  },
  {
    path: "/stocks",
    name: "stocks",
    component: StocksView,
  },
  {
    path: "/recommendations",
    name: "recommendations",
    component: RecommendationsView,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
