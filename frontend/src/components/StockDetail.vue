<template>
  <div class="container mx-auto p-4">
    <div v-if="loading" class="text-center py-4">
      <div class="spinner"></div>
      <p>Cargando detalles del stock...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
      {{ error }}
    </div>

    <div v-else-if="!stock" class="text-center py-8">
      <p class="text-gray-500">Stock no encontrado</p>
      <router-link to="/stocks" class="text-blue-500 hover:underline mt-2 inline-block">
        Volver a la lista de stocks
      </router-link>
    </div>

    <div v-else class="bg-white rounded-lg shadow-lg p-6">
      <div class="flex justify-between items-center mb-6">
        <div>
          <h1 class="text-3xl font-bold">{{ stock.ticker }}</h1>
          <p class="text-xl text-gray-700">{{ stock.company }}</p>
        </div>
        <router-link to="/stocks" class="text-blue-500 hover:underline">
          Volver a la lista
        </router-link>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div class="border-r pr-8">
          <div class="mb-4">
            <span class="text-gray-600">Score:</span>
            <div class="text-3xl font-bold" :class="getScoreColorClass(stock.score)">
              {{ stock.score.toFixed(2) }}
            </div>
          </div>

          <div class="mb-4">
            <span class="text-gray-600">Rating:</span>
            <div class="text-xl font-semibold px-3 py-1 rounded-full inline-block"
              :class="getRatingColorClass(stock.latest_rating)">
              {{ stock.latest_rating }}
            </div>
          </div>

          <div class="mb-4">
            <span class="text-gray-600">Precio objetivo:</span>
            <div class="text-xl font-semibold text-green-600">
              {{ stock.latest_target }}
            </div>
          </div>
        </div>

        <div>
          <h2 class="text-xl font-bold mb-4">Análisis</h2>
          <p class="text-gray-700">{{ stock.analysis_rationale }}</p>
          <p class="text-gray-500 mt-4 text-sm">
            Última actualización: {{ formatDate(stock.last_updated) }}
          </p>
        </div>
      </div>

      <!-- Aquí podrías agregar un gráfico o historial de precios -->
      <div class="mt-8">
        <h2 class="text-xl font-bold mb-4">Historial de precios</h2>
        <div class="h-64 bg-gray-100 rounded flex items-center justify-center">
          <p class="text-gray-500">
            Gráfico de precios (implementación futura)
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useStockStore } from "@/store/stock";
import type { Stock } from "@/types";

export default defineComponent({
  name: "StockDetail",

  setup() {
    const route = useRoute();
    const stockStore = useStockStore();
    const stockTicker = ref(route.params.id as string);

    const loading = computed(() => stockStore.loading);
    const error = computed(() => stockStore.error);

    const stock = computed(() => {
      return stockStore.getStockByTicker(stockTicker.value);
    });

    const getScoreColorClass = (score: number) => {
      if (score >= 70) return "text-green-600";
      if (score >= 50) return "text-yellow-600";
      return "text-red-600";
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

    onMounted(async () => {
      if (!stock.value && !loading.value) {
        await stockStore.fetchStocks();
      }
    });

    watch(
      () => route.params.id,
      (newId) => {
        stockTicker.value = newId as string;
      }
    );

    return {
      stock,
      loading,
      error,
      getScoreColorClass,
      getRatingColorClass,
      formatDate,
    };
  },
});
</script>

<style scoped>
.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border-left-color: #09f;
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
