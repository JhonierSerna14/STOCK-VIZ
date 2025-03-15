import axios from "axios";
import type { Stock, StockResponse, Recommendation } from "@/types";

const apiClient = axios.create({
  baseURL: "http://localhost:8080/api",
  headers: {
    "Content-Type": "application/json",
  },
});

export default {
  async getStocks(page = 1): Promise<StockResponse> {
    const response = await apiClient.get("/stocks/all", {
      params: { page },
    });

    console.log("Respuesta del API:", response.data);

    // Si la respuesta ya viene en el formato esperado
    return response.data;
  },

  async getAllStocks(page = 1): Promise<Stock[]> {
    const response = await apiClient.get("/stocks/all", {
      params: { page },
    });

    console.log("Respuesta getAllStocks:", response.data);

    // Si la respuesta tiene el nuevo formato con items
    if (response.data && response.data.items) {
      return response.data.items;
    }

    // Fallback para compatibilidad
    if (Array.isArray(response.data)) {
      return response.data;
    }

    return [];
  },

  async deleteAllStocks(): Promise<{ mensaje: string }> {
    const response = await apiClient.delete("/stocks");
    return response.data;
  },

  async getRecommendations(filters = {}): Promise<Recommendation[]> {
    const response = await apiClient.get("/recommendations", {
      params: filters,
    });
    return response.data;
  },
};
