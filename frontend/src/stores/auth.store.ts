import { User } from '@/types/user.types';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {

    account: {user: User; accessToken: string; isAuthenticated: boolean } | null;
    _hasHydrated: boolean;
    setHasHydrated: (hydrated: boolean) => void;
    setAuth: (user: User, token: string) => void;
    refreshAccessToken: (user: User, token: string) => void;
    refreshUser: (user: User) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        ((set): AuthState => ({
            account: null,
            _hasHydrated: false,
            setHasHydrated: (hydrated) => set({_hasHydrated: hydrated}),
            setAuth: (user, token) => set({ account: {user, accessToken: token, isAuthenticated: true} }),
            refreshAccessToken: (user, token) => set((state) => {
                if(state.account?.user.id === user.id) {
                    return {
                        account: { ...state.account, isAuthenticated: true, accessToken: token}
                    }
                }
                return state;
            }),
            logout: () => set({ account: null }),
            refreshUser: (user) => set((state) => {
                if(state.account?.user.id === user.id) {
                    return {
                        account: {...state.account, user}
                    }
                }
                return state;
            })
        })),
        {
            name: 'auth-storage',
            partialize: (state) => ({ account: state.account ? {
                user: state.account.user,
                isAuthenticated: false,
                accessToken: null,
            } : null }),
            onRehydrateStorage: () => (state) => {
                if (state) {
                    state.setHasHydrated(true);
                }
            }
        }
    )
)