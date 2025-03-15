<template>
  <div>
    <div class="bg-gray-100 py-6">
      <div class="container mx-auto px-4">
        <h1 class="text-3xl font-bold">Stocks</h1>
        <p class="text-gray-600 mt-2">
          Visualiza información actualizada de stocks del mercado
        </p>
      </div>
    </div>

    <StockList />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import StockList from "@/components/StockList.vue";

export default defineComponent({
  name: "StocksView",
  components: {
    StockList,
  },
});
</script>
