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
    return response.data;
  },

  async getRecommendations(filters = {}): Promise<Recommendation[]> {
    const response = await apiClient.get("/recommendations", {
      params: filters,
    });
    return response.data;
  },
};
