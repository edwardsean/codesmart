export const refreshState = {
    isRefreshing: false,
    pendingQueue: [] as { resolve: (token: string) => void; reject: (err: unknown) => void}[],
};

export const processQueue = (err: unknown, token: string | null) => {
    refreshState.pendingQueue.forEach((p) => err ? p.reject(err) : p.resolve(token!));
    refreshState.pendingQueue = [];
}