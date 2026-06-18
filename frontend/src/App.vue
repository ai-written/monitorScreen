<template>
  <div class="app" @keydown.esc="onEsc" tabindex="0">
    <TopBar :system="data.system" :cpu="data.cpu" :gpu="data.gpu" :memory="data.memory" :total-power="data.total_power" />
    <div class="main-grid">
      <CpuBlock :cpu="data.cpu" />
      <GpuBlock :gpu="data.gpu" />
    </div>
    <div class="bottom-grid">
      <MemoryBlock :memory="data.memory" />
      <StorageBlock :storage="data.storage" />
      <FansBlock :fans="data.fans" />
    </div>
    <div class="overlay-btns">
      <div class="overlay-btn" @click="toggleFs" :title="fullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'">
        <svg v-if="fullscreen" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M8 3v3a2 2 0 01-2 2H3m18 0h-3a2 2 0 01-2-2V3m0 18v-3a2 2 0 012-2h3M3 16h3a2 2 0 012 2v3"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M8 3H5a2 2 0 00-2 2v3m18 0V5a2 2 0 00-2-2h-3m0 18h3a2 2 0 002-2v-3M3 16v3a2 2 0 002 2h3"/>
        </svg>
      </div>
      <div class="overlay-btn" @click="doShutdown" title="Shutdown">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6L6 18M6 6l12 12"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, ref, onMounted, onUnmounted } from 'vue'
import TopBar from './components/TopBar.vue'
import CpuBlock from './components/CpuBlock.vue'
import GpuBlock from './components/GpuBlock.vue'
import MemoryBlock from './components/MemoryBlock.vue'
import StorageBlock from './components/StorageBlock.vue'
import FansBlock from './components/FansBlock.vue'

export default {
  components: { TopBar, CpuBlock, GpuBlock, MemoryBlock, StorageBlock, FansBlock },
  setup() {
    const data = reactive({
      system: { time: '--:--:--', date: '----/--/--', uptime: '--' },
      cpu: { model: '', max_freq: '', cores_threads: '', package_temp: 0, usage: 0, clock_speed: 0, vcore: 0, fan_speed: 0, power: 0 },
      gpu: { model: '', vram_spec: '', pcie_version: '', temp: 0, usage: 0, clock: 0, mem_used: 0, mem_total: 0, fan_speed: 0, power: 0 },
      memory: { used_gb: 0, total_gb: 0, type: '', frequency: '', channel: '', brand: '' },
      storage: [],
      fans: [],
      total_power: 0
    })

    const fullscreen = ref(true)
    let eventCleanup = null
    let clockTimer = null

    function updateClock() {
      const now = new Date()
      const h = String(now.getHours()).padStart(2, '0')
      const m = String(now.getMinutes()).padStart(2, '0')
      const s = String(now.getSeconds()).padStart(2, '0')
      data.system.time = `${h}:${m}:${s}`
      const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
      const y = now.getFullYear()
      const mo = String(now.getMonth() + 1).padStart(2, '0')
      const d = String(now.getDate()).padStart(2, '0')
      data.system.date = `${y}/${mo}/${d} ${weekdays[now.getDay()]}`
    }

    onMounted(async () => {
      updateClock()
      clockTimer = setInterval(updateClock, 1000)

      try {
        const d = await window.go.main.App.GetDashboard()
        if (d) {
          if (d.system) {
            data.system.uptime = d.system.uptime
            delete d.system
          }
          Object.assign(data, d)
        }
        fullscreen.value = await window.go.main.App.GetIsFullscreen()
      } catch (e) { console.error(e) }

      if (window.runtime) {
        eventCleanup = window.runtime.EventsOn('dashboard-update', (raw) => {
          try {
            const d = typeof raw === 'string' ? JSON.parse(raw) : raw
            if (d) {
              if (d.system) {
                data.system.uptime = d.system.uptime
                delete d.system
              }
              Object.assign(data, d)
            }
          } catch (e) { /* ignore */ }
        })
      }

      document.addEventListener('keydown', handleKey)
    })

    onUnmounted(() => {
      if (clockTimer) clearInterval(clockTimer)
      if (eventCleanup) eventCleanup()
      document.removeEventListener('keydown', handleKey)
    })

    function handleKey(e) {
      if (e.key === 'Escape') {
        toggleFs()
      }
    }

    async function toggleFs() {
      try {
        await window.go.main.App.ToggleFullscreen()
        fullscreen.value = await window.go.main.App.GetIsFullscreen()
      } catch (e) { /* ignore */ }
    }

    async function doShutdown() {
      try { await window.go.main.App.Shutdown() } catch (e) { /* ignore */ }
    }

    function onEsc() {}

    return { data, fullscreen, toggleFs, doShutdown, onEsc }
  }
}
</script>

<style scoped>
.app {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 12px 16px;
  gap: 10px;
}

.main-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  flex: 1;
  min-height: 0;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
  flex-shrink: 0;
}

.overlay-btns {
  position: fixed;
  top: 6px;
  right: 6px;
  display: flex;
  gap: 4px;
  z-index: 100;
}

.overlay-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  cursor: pointer;
  opacity: 0.15;
  transition: opacity 0.2s;
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
  -webkit-app-region: no-drag;
}

.overlay-btn:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.2);
}
</style>
