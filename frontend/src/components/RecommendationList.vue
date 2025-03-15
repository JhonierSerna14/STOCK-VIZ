<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6 text-gray-800">
      Recomendaciones de inversión
    </h1>

    <!-- Panel de filtros -->
    <div class="bg-white p-6 rounded-lg shadow-md mb-8 border border-gray-100">
      <h2 class="text-lg font-semibold mb-4 flex items-center">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 mr-2 text-blue-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
          />
        </svg>
        Filtros
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
        <!-- Rating -->
        <div>
          <label
            for="rating"
            class="block text-sm font-medium text-gray-700 mb-1"
            >Rating</label
          >
          <select
            id="rating"
            v-model="filters.rating"
            class="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-300 focus:border-blue-500 focus:outline-none transition duration-200"
          >
            <option value="">Todos</option>
            <option value="buy">Buy</option>
            <option value="overweight">Overweight</option>
            <option value="hold">Hold</option>
            <option value="neutral">Neutral</option>
            <option value="sell">Sell</option>
            <option value="underweight">Underweight</option>
          </select>
        </div>

        <!-- Límite -->
        <div>
          <label
            for="limit"
            class="block text-sm font-medium text-gray-700 mb-1"
            >Cantidad</label
          >
          <input
            type="number"
            id="limit"
            v-model.number="filters.limit"
            min="1"
            max="50"
            class="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-300 focus:border-blue-500 focus:outline-none transition duration-200"
          />
        </div>

        <!-- Fecha desde -->
        <div>
          <label
            for="date_from"
            class="block text-sm font-medium text-gray-700 mb-1"
            >Desde</label
          >
          <input
            type="date"
            id="date_from"
            v-model="filters.date_from"
            class="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-300 focus:border-blue-500 focus:outline-none transition duration-200"
          />
        </div>

        <!-- Fecha hasta -->
        <div>
          <label
            for="date_to"
            class="block text-sm font-medium text-gray-700 mb-1"
            >Hasta</label
          >
          <input
            type="date"
            id="date_to"
            v-model="filters.date_to"
            class="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-300 focus:border-blue-500 focus:outline-none transition duration-200"
          />
        </div>
      </div>

      <div class="flex space-x-3">
        <button
          @click="applyFilters"
          class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md transition duration-300 ease-in-out flex items-center shadow-sm"
          :disabled="loading"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5 mr-1"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 5l7 7-7 7"
            />
          </svg>
          Aplicar filtros
        </button>
        <button
          @click="clearFilters"
          class="bg-gray-200 hover:bg-gray-300 text-gray-700 px-4 py-2 rounded-md transition duration-300 ease-in-out flex items-center shadow-sm"
          :disabled="loading"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5 mr-1"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
          Limpiar filtros
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8 bg-white rounded-lg shadow-md">
      <div class="spinner"></div>
      <p class="mt-4 text-gray-600">Cargando recomendaciones...</p>
    </div>

    <div
      v-else-if="error"
      class="bg-red-100 border border-red-400 text-red-700 px-6 py-4 rounded-lg shadow-sm my-4"
    >
      <div class="flex items-center">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6 mr-2 text-red-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
        {{ error }}
      </div>
    </div>

    <div
      v-else-if="recommendations.length === 0"
      class="text-center py-12 bg-white rounded-lg shadow-md"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-16 w-16 mx-auto text-gray-400 mb-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <p class="text-gray-500 text-lg">No hay recomendaciones disponibles</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
      <div
        v-for="recommendation in recommendations"
        :key="recommendation.ticker"
        class="border rounded-xl p-6 shadow-sm hover:shadow-lg transition-all duration-300 bg-white overflow-hidden transform hover:-translate-y-1"
      >
        <!-- Cabecera con ticker y nombre de empresa -->
        <div
          class="flex justify-between items-start mb-4 pb-3 border-b border-gray-100"
        >
          <div>
            <div class="flex items-center">
              <h3 class="text-2xl font-bold text-blue-600 mr-2">
                {{ recommendation.ticker }}
              </h3>
              <span
                class="text-xs px-2 py-1 rounded-full uppercase font-bold tracking-wider"
                :class="getRatingColorClass(recommendation.latest_rating)"
              >
                {{ recommendation.latest_rating }}
              </span>
            </div>
            <p class="text-gray-600 text-lg">{{ recommendation.company }}</p>
          </div>

          <!-- Score con indicador visual -->
          <div class="bg-gray-50 rounded-lg p-3 shadow-inner">
            <div class="text-center">
              <div class="font-semibold text-sm text-gray-500 mb-1">Score</div>
              <div
                class="text-2xl font-bold rounded-full w-12 h-12 flex items-center justify-center mx-auto"
                :class="getScoreBackgroundClass(recommendation.score)"
              >
                {{ Math.round(recommendation.score) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Contenido principal -->
        <div class="mb-4">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500 mr-1"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                />
              </svg>
              <span class="font-medium text-gray-700 mr-1">Target: </span>
              <span class="text-green-600 font-bold">
                {{ recommendation.latest_target }}
              </span>
            </div>

            <div class="text-gray-500 text-xs flex items-center">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-4 w-4 mr-1"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              {{ formatDate(recommendation.last_updated) }}
            </div>
          </div>

          <div class="bg-gray-50 p-4 rounded-lg border border-gray-100">
            <h4 class="font-medium text-gray-700 mb-2 flex items-center">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-4 w-4 mr-1 text-blue-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
              Análisis
            </h4>
            <p class="text-gray-700">
              {{ recommendation.analysis_rationale }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, computed, ref } from "vue";
import { useRecommendationStore } from "@/store/recommendation";
import type { RecommendationFilter } from "@/types";

export default defineComponent({
  name: "RecommendationList",

  setup() {
    const recommendationStore = useRecommendationStore();

    // Estado local para los filtros
    const filters = ref<RecommendationFilter>({
      limit: 5,
      date_from: "",
      date_to: "",
      rating: "",
      ticker: "",
    });

    onMounted(() => {
      recommendationStore.fetchRecommendations();
    });

    const refreshRecommendations = () => {
      recommendationStore.fetchRecommendations();
    };

    const applyFilters = () => {
      recommendationStore.fetchRecommendations(filters.value);
    };

    const clearFilters = () => {
      filters.value = {
        limit: 5,
        date_from: "",
        date_to: "",
        rating: "",
        ticker: "",
      };
      recommendationStore.clearFilters();
    };

    const getScoreColorClass = (score: number) => {
      if (score >= 70) return "text-green-600";
      if (score >= 50) return "text-yellow-600";
      return "text-red-600";
    };

    const getScoreBackgroundClass = (score: number) => {
      if (score >= 70) return "bg-green-100 text-green-800";
      if (score >= 50) return "bg-yellow-100 text-yellow-800";
      return "bg-red-100 text-red-800";
    };

    const getRatingColorClass = (rating: string) => {
      const ratingLower = rating.toLowerCase();
      if (ratingLower.includes("buy") || ratingLower.includes("overweight")) {
        return "bg-green-100 text-green-800";
      }
      if (ratingLower.includes("hold") || ratingLower.includes("neutral")) {
        return "bg-yellow-100 text-yellow-800";
      }
      return "bg-red-100 text-red-800";
    };

    const formatDate = (dateString: string) => {
      return new Date(dateString).toLocaleDateString("es-ES", {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    };

    return {
      recommendations: computed(() => recommendationStore.recommendations),
      loading: computed(() => recommendationStore.loading),
      error: computed(() => recommendationStore.error),
      filters,
      refreshRecommendations,
      applyFilters,
      clearFilters,
      getScoreColorClass,
      getScoreBackgroundClass,
      getRatingColorClass,
      formatDate,
    };
  },
});
</script>

<style scoped>
.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-left-color: #3498db;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
