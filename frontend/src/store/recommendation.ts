import { defineStore } from "pinia";
import type { Recommendation, RecommendationFilter } from "@/types";
import api from "@/services/api";

export const useRecommendationStore = defineStore("recommendation", {
  state: () => ({
    recommendations: [] as Recommendation[],
    loading: false,
    error: null as string | null,
    filters: {
      limit: 6,
      date_from: "",
      date_to: "",
      rating: "",
      ticker: "",
    } as RecommendationFilter,
  }),

  actions: {
    async fetchRecommendations(newFilters?: Partial<RecommendationFilter>) {
      try {
        // Si se proporcionan nuevos filtros, actualizar los filtros existentes
        if (newFilters) {
          this.filters = {
            ...this.filters,
            ...newFilters,
          };
        }

        this.loading = true;
        const recommendations = await api.getRecommendations(this.filters);
        this.recommendations = recommendations;
        this.error = null;
      } catch (error) {
        this.error = "Error al cargar las recomendaciones";
        console.error(error);
      } finally {
        this.loading = false;
      }
    },

    // Acción para limpiar todos los filtros
    clearFilters() {
      this.filters = {
        limit: 6,
        date_from: "",
        date_to: "",
        rating: "",
        ticker: "",
      };
      this.fetchRecommendations();
    },
  },
});
