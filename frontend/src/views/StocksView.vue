<template>
  <div>
    <div class="bg-gray-100 py-6">
      <div class="container mx-auto px-4">
        <h1 class="text-3xl font-bold">Stocks</h1>
        <div class="flex flex-col md:flex-row md:items-center md:justify-between mt-2">
          <p class="text-gray-600">
            Visualiza información actualizada de stocks del mercado
          </p>
          <div class="mt-4 md:mt-0 md:ml-4 w-full md:w-auto max-w-sm">
            <div class="flex overflow-hidden rounded-md bg-gray-200 focus:outline focus:outline-blue-500">
              <input type="text" placeholder="Buscar" v-model="searchQuery"
                class="w-full rounded-bl-md rounded-tl-md bg-gray-100 px-4 py-2.5 text-gray-700 focus:outline-blue-500" />
              <button @click="searchStocks" class="bg-blue-500 px-3.5 text-white duration-150 hover:bg-blue-600">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="size-6">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <StockList :search-query="searchQuery" @search="searchStocks" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import StockList from "@/components/StockList.vue";
import { useStockStore } from "@/store/stock";


export default defineComponent({
  name: "StocksView",
  components: {
    StockList,
  },
  setup() {
    const stockStore = useStockStore();
    const searchQuery = ref("");

    const searchStocks = async () => {
      await stockStore.fetchStocks(1, searchQuery.value);
    };

    // Inicializar con el término de búsqueda guardado en el store
    searchQuery.value = stockStore.currentSearchQuery;

    return {
      searchStocks,
      searchQuery,
    };
  },
});
</script>
