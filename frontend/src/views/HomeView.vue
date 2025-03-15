<template>
  <div class="container mx-auto p-4">
    <div class="py-12 text-center">
      <h1 class="text-4xl font-bold mb-4">Bienvenido a StockViz</h1>
      <p class="text-xl text-gray-600 max-w-2xl mx-auto mb-8">
        Tu plataforma para visualizar datos de acciones y obtener
        recomendaciones de inversión basadas en análisis avanzados.
      </p>

      <div class="flex flex-col md:flex-row justify-center gap-4 mt-8">
        <router-link
          to="/stocks"
          class="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium transition-colors"
        >
          Ver Stocks
        </router-link>

        <router-link
          to="/recommendations"
          class="bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-lg font-medium transition-colors"
        >
          Ver Recomendaciones
        </router-link>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8 my-12">
      <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="text-blue-500 mb-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-10 w-10"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"
            />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-2">Visualización de Stocks</h3>
        <p class="text-gray-600">
          Accede a información actualizada sobre stocks y visualiza su
          rendimiento de manera intuitiva.
        </p>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="text-green-500 mb-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-10 w-10"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-2">Análisis Inteligente</h3>
        <p class="text-gray-600">
          Nuestro sistema analiza los datos para identificar patrones y
          tendencias que pueden influir en tus decisiones.
        </p>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="text-purple-500 mb-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-10 w-10"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z"
            />
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z"
            />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-2">Recomendaciones Personalizadas</h3>
        <p class="text-gray-600">
          Obtén recomendaciones de inversión basadas en análisis de datos y
          algoritmos avanzados.
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "HomeView",
});
</script>
