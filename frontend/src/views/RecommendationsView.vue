<template>
  <div>
    <div class="bg-gray-100 py-6">
      <div class="container mx-auto px-4">
        <h1 class="text-3xl font-bold">Recomendaciones</h1>
        <p class="text-gray-600 mt-2">
          Recibe recomendaciones de inversión basadas en análisis de datos
        </p>
      </div>
    </div>

    <RecommendationList />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import RecommendationList from "@/components/RecommendationList.vue";

export default defineComponent({
  name: "RecommendationsView",
  components: {
    RecommendationList,
  },
});
</script>
