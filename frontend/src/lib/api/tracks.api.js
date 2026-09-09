import apiClient from './client.js';

export const tracksApi = {
  getAll: async () => {
    return apiClient.get('/tracks');
  },

  getById: async (id) => {
    return apiClient.get(`/tracks/${id}`);
  },

  create: async (formData) => {
    return apiClient.post('/tracks', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  addMoodItem: async (trackId, formData) => {
    return apiClient.post(`/tracks/${trackId}/items`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  extractStream: async (payload) => {
    return apiClient.post('/tracks/extract', payload);
  },

  getExtractJob: async (jobId) => {
    return apiClient.get(`/tracks/extract/${jobId}`);
  },

  delete: async (id) => {
    return apiClient.delete(`/tracks/${id}`);
  },
};
