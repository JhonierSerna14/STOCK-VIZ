<template>
  <header class="bg-gray-800 text-white shadow-md">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between py-4">
        <div class="flex items-center">
          <router-link to="/" class="text-xl font-bold">StockViz</router-link>

          <nav class="ml-10 hidden md:flex space-x-6">
            <router-link
              to="/"
              class="hover:text-blue-300 transition-colors"
              :class="{ 'text-blue-300': currentRoute === '/' }"
            >
              Inicio
            </router-link>
            <router-link
              to="/stocks"
              class="hover:text-blue-300 transition-colors"
              :class="{ 'text-blue-300': currentRoute === '/stocks' }"
            >
              Stocks
            </router-link>
            <router-link
              to="/recommendations"
              class="hover:text-blue-300 transition-colors"
              :class="{ 'text-blue-300': currentRoute === '/recommendations' }"
            >
              Recomendaciones
            </router-link>
          </nav>
        </div>

        <!-- Botón de menú móvil -->
        <div class="md:hidden">
          <button @click="toggleMobileMenu" class="focus:outline-none">
            <svg class="h-6 w-6 fill-current" viewBox="0 0 24 24">
              <path
                v-if="!mobileMenuOpen"
                d="M4 6h16v2H4zm0 5h16v2H4zm0 5h16v2H4z"
              ></path>
              <path
                v-else
                d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
              ></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- Menú móvil -->
      <div v-show="mobileMenuOpen" class="md:hidden pb-4">
        <nav class="flex flex-col space-y-3">
          <router-link
            to="/"
            @click="closeMobileMenu"
            class="hover:text-blue-300 transition-colors py-1"
            :class="{ 'text-blue-300': currentRoute === '/' }"
          >
            Inicio
          </router-link>
          <router-link
            to="/stocks"
            @click="closeMobileMenu"
            class="hover:text-blue-300 transition-colors py-1"
            :class="{ 'text-blue-300': currentRoute === '/stocks' }"
          >
            Stocks
          </router-link>
          <router-link
            to="/recommendations"
            @click="closeMobileMenu"
            class="hover:text-blue-300 transition-colors py-1"
            :class="{ 'text-blue-300': currentRoute === '/recommendations' }"
          >
            Recomendaciones
          </router-link>
        </nav>
      </div>
    </div>
  </header>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import { useRoute } from "vue-router";

export default defineComponent({
  name: "AppHeader",

  setup() {
    const route = useRoute();
    const mobileMenuOpen = ref(false);

    const currentRoute = computed(() => route.path);

    const toggleMobileMenu = () => {
      mobileMenuOpen.value = !mobileMenuOpen.value;
    };

    const closeMobileMenu = () => {
      mobileMenuOpen.value = false;
    };

    return {
      mobileMenuOpen,
      currentRoute,
      toggleMobileMenu,
      closeMobileMenu,
    };
  },
});
</script>
