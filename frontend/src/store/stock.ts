import { defineStore } from "pinia";
import type { Stock, StockResponse, Pagination } from "@/types";
import api from "@/services/api";

export const useStockStore = defineStore("stock", {
  state: () => ({
    stocks: [] as Stock[],
    pagination: {
      current_page: 1,
      per_page: 20,
      total_items: 0,
      total_pages: 0,
    } as Pagination,
    loading: false,
    error: null as string | null,
    currentSearchQuery: "",
  }),

  getters: {
    getStockByTicker: (state) => (ticker: string) => {
      return state.stocks.find((stock) => stock.ticker === ticker);
    },
    currentPage: (state) => state.pagination.current_page,
    totalPages: (state) => state.pagination.total_pages,
    hasNextPage: (state) =>
      state.pagination.current_page < state.pagination.total_pages,
    hasPrevPage: (state) => state.pagination.current_page > 1,
  },

  actions: {
    async fetchStocks(page = 1, search = "") {
      try {
        this.loading = true;
        this.currentSearchQuery = search;
        const response = await api.getStocks(page, search);

        if (response.items) {
          this.stocks = response.items;
          this.pagination = response.pagination;
        } else {
          console.error("Formato de respuesta no esperado");
        }

        this.error = null;
      } catch (error) {
        console.error("Error detallado en fetchStocks:", error);
        this.error = "Error al cargar los stocks";
      } finally {
        this.loading = false;
      }
    },

    async goToPage(page: number) {
      if (page >= 1 && page <= this.pagination.total_pages) {
        await this.fetchStocks(page, this.currentSearchQuery);
      }
    },

    async nextPage() {
      if (this.pagination.current_page < this.pagination.total_pages) {
        await this.fetchStocks(this.pagination.current_page + 1, this.currentSearchQuery);
      }
    },

    async prevPage() {
      if (this.pagination.current_page > 1) {
        await this.fetchStocks(this.pagination.current_page - 1, this.currentSearchQuery);
      }
    },
  },
});
