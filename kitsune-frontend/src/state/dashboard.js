import { create } from 'zustand'

export const useDashboardState = create((set, get) => ({
    selectedImplants: [],
    newTaskWindowOpen: false,
    showCompletedTasks: false,
    selectImplant: (implantId) => {
      if(get().selectedImplants.includes(implantId)){
        set({selectedImplants: get().selectedImplants.filter((i) => (i != implantId))})
      } else {
        set({selectedImplants: [...get().selectedImplants, implantId]})
      }
    },
    setNewTaskWindowOpen: (val) => set({newTaskWindowOpen: val}),
    setShowCompletedTasks: (val) => set({showCompletedTasks: val}),
  }))