export const queryKeys = {
  tracks: {
    all: ['tracks'],
    detail: (id) => ['tracks', id],
  },
  auth: {
    me: ['auth', 'me'],
  },
};
