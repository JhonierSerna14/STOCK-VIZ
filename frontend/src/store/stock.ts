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
    async fetchStocks(page = 1) {
      try {
        this.loading = true;
        console.log(`Llamando a api.getStocks para página ${page}...`);
        const response = await api.getStocks(page);
        console.log("Respuesta recibida:", response);

        if (response.items) {
          this.stocks = response.items;
          this.pagination = response.pagination;
        } else {
          console.log("Formato de respuesta no esperado");
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
        await this.fetchStocks(page);
      }
    },

    async nextPage() {
      if (this.pagination.current_page < this.pagination.total_pages) {
        await this.fetchStocks(this.pagination.current_page + 1);
      }
    },

    async prevPage() {
      if (this.pagination.current_page > 1) {
        await this.fetchStocks(this.pagination.current_page - 1);
      }
    },

    async fetchAllStocks(page = 1) {
      try {
        this.loading = true;
        console.log(`Llamando a api.getAllStocks para página ${page}...`);
        const stocks = await api.getAllStocks(page);
        console.log("Stocks recibidos en fetchAllStocks:", stocks.length);
        this.stocks = stocks;
        this.error = null;
      } catch (error) {
        console.error("Error detallado en fetchAllStocks:", error);
        this.error = "Error al cargar todos los stocks";
      } finally {
        this.loading = false;
      }
    },

    async deleteAllStocks() {
      try {
        this.loading = true;
        await api.deleteAllStocks();
        this.stocks = [];
        this.error = null;
      } catch (error) {
        this.error = "Error al eliminar los stocks";
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
  },
});
