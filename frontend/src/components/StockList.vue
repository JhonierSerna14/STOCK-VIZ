<template>
  <div class="container mx-auto p-4">
    <div v-if="loading" class="text-center py-4">
      <div class="spinner"></div>
      <p>Cargando actualizaciones de stocks...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
      {{ error }}
    </div>

    <div v-else-if="stocks?.length === 0" class="text-center py-8">
      <p class="text-gray-500">No hay actualizaciones disponibles</p>
    </div>

    <div v-else>
      <div class="grid grid-cols-1 gap-4">
        <div v-for="stock in stocks" :key="`${stock.ticker}-${stock.time}`"
          class="border rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow bg-white">
          <div class="flex justify-between items-start mb-4">
            <div>
              <h3 class="text-2xl font-bold text-blue-600">
                {{ stock.ticker }}
              </h3>
              <p class="text-gray-600 text-lg">{{ stock.company }}</p>
            </div>
            <div class="text-right">
              <span class="px-3 py-1 rounded-full text-sm font-medium" :class="getActionColorClass(stock.action)">
                {{ formatAction(stock.action) }}
              </span>
            </div>
          </div>

          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="mb-4">
              <p class="font-medium text-gray-700">Por {{ stock.brokerage }}</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-600 mb-1">Rating:</p>
                <div class="flex items-center">
                  <span class="px-2 py-1 rounded bg-gray-100 text-sm" :class="getRatingColorClass(stock.rating_from)">
                    {{ stock.rating_from }}
                  </span>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mx-2 text-gray-400" fill="none"
                    viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                  <span class="px-2 py-1 rounded text-sm font-medium" :class="getRatingColorClass(stock.rating_to)">
                    {{ stock.rating_to }}
                  </span>
                </div>
              </div>

              <div>
                <p class="text-sm text-gray-600 mb-1">Precio objetivo:</p>
                <div class="flex items-center">
                  <span class="px-2 py-1 rounded bg-gray-100 text-sm">
                    {{ stock.target_from }}
                  </span>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mx-2 text-gray-400" fill="none"
                    viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                  <span class="px-2 py-1 rounded text-sm font-medium" :class="getTargetChangeClass(stock.target_from, stock.target_to)
                    ">
                    {{ stock.target_to }}
                  </span>
                </div>
              </div>
            </div>

            <p class="text-gray-500 text-xs mt-4">
              {{ formatDate(stock.time) }}
            </p>
          </div>
        </div>
      </div>

      <div class="flex justify-between items-center mt-8">
        <div class="text-sm text-gray-600">
          Mostrando {{ stocks.length }} de {{ totalItems }} resultados
        </div>
        <div class="flex space-x-2">
          <button @click="prevPage" class="px-4 py-2 bg-gray-200 rounded-md disabled:opacity-50"
            :disabled="!hasPrevPage" :class="{ 'hover:bg-gray-300': hasPrevPage }">
            Anterior
          </button>

          <div class="flex items-center px-4">
            <span>{{ currentPage }} de {{ totalPages }}</span>
          </div>

          <button @click="nextPage" class="px-4 py-2 bg-gray-200 rounded-md disabled:opacity-50"
            :disabled="!hasNextPage" :class="{ 'hover:bg-gray-300': hasNextPage }">
            Siguiente
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, computed } from "vue";
import { useStockStore } from "@/store/stock";

export default defineComponent({
  name: "StockList",

  setup() {
    const stockStore = useStockStore();

    onMounted(async () => {
      stockStore.fetchStocks();
    });

    const getActionColorClass = (action: string) => {
      const actionLower = action.toLowerCase();
      if (actionLower.includes("upgraded") || actionLower.includes("raised")) {
        return "bg-green-100 text-green-800";
      }
      if (actionLower.includes("downgraded") || actionLower.includes("lowered")) {
        return "bg-red-100 text-red-800";
      }
      if (actionLower.includes("reiterated")) {
        return "bg-blue-100 text-blue-800";
      }
      if (actionLower.includes("initiated")) {
        return "bg-purple-100 text-purple-800";
      }
      return "bg-gray-100 text-gray-800";
    };

    const formatAction = (action: string) => {
      const actionMap: { [key: string]: string } = {
        "upgraded by": "Mejora",
        "downgraded by": "Rebaja",
        "reiterated by": "Reiteración",
        "initiated by": "Inicio",
        "target raised by": "Objetivo aumentado",
        "target lowered by": "Objetivo reducido",
        "target set by": "Objetivo establecido",
      };
      return actionMap[action] || action;
    };

    const getRatingColorClass = (rating: string) => {
      const ratingLower = rating.toLowerCase();
      if (
        ratingLower.includes("buy") ||
        ratingLower.includes("overweight") ||
        ratingLower.includes("outperform")
      ) {
        return "bg-green-100 text-green-800";
      }
      if (
        ratingLower.includes("hold") ||
        ratingLower.includes("neutral") ||
        ratingLower.includes("market perform") ||
        ratingLower.includes("equal weight")
      ) {
        return "bg-yellow-100 text-yellow-800";
      }
      if (
        ratingLower.includes("sell") ||
        ratingLower.includes("underweight") ||
        ratingLower.includes("underperform")
      ) {
        return "bg-red-100 text-red-800";
      }
      return "bg-gray-100 text-gray-800";
    };

    const getTargetChangeClass = (from: string, to: string) => {
      const fromValue = parseFloat(from.replace("$", ""));
      const toValue = parseFloat(to.replace("$", ""));
      if (toValue > fromValue) {
        return "bg-green-100 text-green-800";
      }
      if (toValue < fromValue) {
        return "bg-red-100 text-red-800";
      }
      return "bg-gray-100 text-gray-800";
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

    const nextPage = async () => {
      await stockStore.nextPage();
    };

    const prevPage = async () => {
      await stockStore.prevPage();
    };

    return {
      stocks: computed(() => stockStore.stocks),
      loading: computed(() => stockStore.loading),
      error: computed(() => stockStore.error),
      currentPage: computed(() => stockStore.currentPage),
      totalPages: computed(() => stockStore.totalPages),
      totalItems: computed(() => stockStore.pagination.total_items),
      hasNextPage: computed(() => stockStore.hasNextPage),
      hasPrevPage: computed(() => stockStore.hasPrevPage),
      getActionColorClass,
      formatAction,
      getRatingColorClass,
      getTargetChangeClass,
      formatDate,
      nextPage,
      prevPage,
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
