import { create } from 'zustand'

export const useDashboardState = create((set, get) => ({
    selectedImplants: [],
    newTaskWindowOpen: false,
    resultWindowOpen: false,
    taskResult: {},
    showCompletedTasks: false,
    notificationQueue: [],
    pushNotification: (val) => {
      set({notificationQueue: [...get().notificationQueue, val]})
    },
    popNotification: () => {
      const queue = get().notificationQueue
      const poppedItem = queue[queue.length - 1]
      set({notificationQueue: queue.slice(0, -1)})
      return poppedItem
    },
    selectImplant: (implantId) => {
      if(get().selectedImplants.includes(implantId)){
        set({selectedImplants: get().selectedImplants.filter((i) => (i != implantId))})
      } else {
        set({selectedImplants: [...get().selectedImplants, implantId]})
      }
    },
    setNewTaskWindowOpen: (val) => set({newTaskWindowOpen: val}),
    setResultWindowOpen: (val) => set({resultWindowOpen: val}),
    setTaskResult: (val) => set({taskResult: val}),
    setShowCompletedTasks: (val) => set({showCompletedTasks: val}),
  }))