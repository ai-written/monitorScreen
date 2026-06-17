<template>
  <div ref="topbar" class="topbar">
    <div class="topbar-left">
      <div class="clock">{{ system.time }}</div>
      <div class="date">{{ system.date }}</div>
    </div>
    <div class="topbar-center">
      <div class="uptime-label">内核运行时间</div>
      <div class="uptime-value">{{ system.uptime }}</div>
    </div>
    <div class="topbar-right">
      <div class="stat-item">
        <span class="stat-label">CPU</span>
        <span class="stat-value" :class="barColor(cpu.usage)">{{ fmt(cpu.usage) }}%</span>
        <div class="mini-bar"><div class="mini-fill" :style="{ width: cpu.usage + '%', background: barColor(cpu.usage) }"></div></div>
      </div>
      <div class="stat-item">
        <span class="stat-label">GPU</span>
        <span class="stat-value" :class="barColor(gpu.usage)">{{ fmt(gpu.usage) }}%</span>
        <div class="mini-bar"><div class="mini-fill" :style="{ width: gpu.usage + '%', background: barColor(gpu.usage) }"></div></div>
      </div>
      <div class="stat-item">
        <span class="stat-label">RAM</span>
        <span class="stat-value" :class="barColor(memPct)">{{ fmt(memPct) }}%</span>
        <div class="mini-bar"><div class="mini-fill" :style="{ width: memPct + '%', background: barColor(memPct) }"></div></div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  props: { system: Object, cpu: Object, gpu: Object, memory: Object },
  setup(props) {
    const topbar = ref(null)
    let dragging = false
    let onMove = null
    let onUp = null

    const memPct = computed(() => {
      if (props.memory && props.memory.total_gb > 0) {
        return (props.memory.used_gb / props.memory.total_gb) * 100
      }
      return 0
    })

    function fmt(v) { return typeof v === 'number' ? v.toFixed(1) : '--' }
    function barColor(v) {
      if (v > 90) return 'var(--red)'
      if (v > 70) return 'var(--amber)'
      return 'var(--green)'
    }

    function cleanup() {
      if (onMove) document.removeEventListener('mousemove', onMove)
      if (onUp) document.removeEventListener('mouseup', onUp)
      onMove = null
      onUp = null
      document.body.style.cursor = ''
      dragging = false
      try { window.go.main.App.DragEnd() } catch (_) {}
    }

    function onDragStart(e) {
      if (e.button !== 0 || dragging) return
      cleanup()

      dragging = true
      try { window.go.main.App.DragStart(e.screenX, e.screenY) } catch (_) {}

      onMove = function(ev) {
        try { window.go.main.App.DragMove(ev.screenX, ev.screenY) } catch (_) {}
      }

      onUp = function() {
        cleanup()
      }

      document.addEventListener('mousemove', onMove)
      document.addEventListener('mouseup', onUp)
      document.body.style.cursor = 'grabbing'
    }

    onMounted(() => {
      if (topbar.value) {
        topbar.value.addEventListener('mousedown', onDragStart)
      }
    })

    onUnmounted(() => {
      cleanup()
      if (topbar.value) {
        topbar.value.removeEventListener('mousedown', onDragStart)
      }
    })

    return { topbar, memPct, fmt, barColor }
  }
}
</script>

<style scoped>
.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px 24px;
  flex-shrink: 0;
  cursor: grab;
  user-select: none;
  -webkit-user-select: none;
}

.clock { font-size: 32px; font-weight: 700; letter-spacing: 2px; color: var(--cyan); }
.date { font-size: 13px; color: var(--text-secondary); margin-top: 2px; }

.topbar-center { text-align: center; }
.uptime-label { font-size: 10px; color: var(--text-dim); letter-spacing: 2px; margin-bottom: 2px; }
.uptime-value { font-size: 16px; font-weight: 600; color: var(--text-primary); }

.topbar-right { display: flex; gap: 24px; }
.stat-item { text-align: center; }
.stat-label { font-size: 10px; color: var(--text-dim); display: block; letter-spacing: 1px; }
.stat-value { font-size: 18px; font-weight: 700; display: block; }

.mini-bar {
  width: 48px; height: 3px; background: var(--border); border-radius: 2px;
  margin-top: 4px; overflow: hidden;
}
.mini-fill { height: 100%; border-radius: 2px; transition: width 0.3s; }
</style>
