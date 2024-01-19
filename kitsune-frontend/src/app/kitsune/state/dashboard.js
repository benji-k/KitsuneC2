import { create } from 'zustand'

export const useDashboardState = create((set) => ({
    selectedImplant: "",
    newTaskWindowOpen: false,
    showCompletedTasks: false,
    setSelectedImplant: (implantId) => set({selectedImplant: implantId}),
    setNewTaskWindowOpen: (val) => set({newTaskWindowOpen: val}),
    setShowCompletedTasks: (val) => set({showCompletedTasks: val}),
  }))