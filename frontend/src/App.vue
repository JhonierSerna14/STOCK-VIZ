<template>
  <div class="min-h-screen flex flex-col bg-gray-50">
    <AppHeader />

    <main class="flex-grow">
      <router-view />
    </main>

    <footer class="bg-gray-800 text-white py-6">
      <div class="container mx-auto px-4">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="mb-4 md:mb-0">
            <p>
              &copy; {{ currentYear }} StockViz. Todos los derechos reservados.
            </p>
          </div>
          <div class="flex space-x-4">
            <a href="#" class="hover:text-blue-300 transition-colors"
              >Términos</a
            >
            <a href="#" class="hover:text-blue-300 transition-colors"
              >Privacidad</a
            >
            <a href="#" class="hover:text-blue-300 transition-colors"
              >Contacto</a
            >
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from "vue";
import AppHeader from "@/components/AppHeader.vue";

export default defineComponent({
  name: "App",
  components: {
    AppHeader,
  },
  setup() {
    const currentYear = computed(() => new Date().getFullYear());

    return {
      currentYear,
    };
  },
});
</script>
