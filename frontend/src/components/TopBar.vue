<template>
  <div ref="topbar" class="topbar">
    <div class="topbar-left">
      <div class="clock">{{ system.time }}</div>
      <div class="date">{{ system.date }}</div>
    </div>
    <div class="topbar-center">
      <div class="uptime-label">内核运行时间</div>
      <div class="uptime-value">{{ system.uptime }}</div>
      <div class="sys-info">
        <span v-if="system.ip_address">{{ system.ip_address }}</span>
        <span v-if="system.ip_address && system.display_info" class="sys-sep">|</span>
        <span v-if="system.display_info">{{ system.display_info }}</span>
        <span v-if="system.fps > 0" class="sys-sep">|</span>
        <span v-if="system.fps > 0" class="sys-fps">{{ system.fps }} <span style="font-weight:400;opacity:0.8">UI FPS</span></span>
      </div>
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
      <div class="stat-item">
        <span class="stat-label">功耗</span>
        <span class="stat-value pwr-val">{{ pwrDisplay }}<span class="pwr-unit">W</span></span>
        <div class="mini-bar"><div class="mini-fill pwr-fill" :style="{ width: pwrPct + '%' }"></div></div>
      </div>
      <div class="stat-item">
        <span class="stat-label">网络 ↓/↑</span>
        <span class="stat-value net-val">{{ netDown }} / {{ netUp }}<span class="net-unit">Mbps</span></span>
        <div class="mini-bar"><div class="mini-fill net-fill" :style="{ width: netPct + '%' }"></div></div>
      </div>
      <div class="stat-item" v-if="system.process_count > 0">
        <span class="stat-label">进程 / 线程</span>
        <span class="stat-value proc-val">{{ system.process_count }} / {{ system.thread_count || '--' }}</span>
        <div class="mini-bar mini-bar-ghost"><div class="mini-fill"></div></div>
      </div>
      <div class="stat-item" v-if="network.connection_count > 0">
        <span class="stat-label">Tcp 连接</span>
        <span class="stat-value proc-val">{{ network.connection_count }}</span>
        <div class="mini-bar mini-bar-ghost"><div class="mini-fill"></div></div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  props: { system: Object, cpu: Object, gpu: Object, memory: Object, totalPower: { type: Number, default: 0 }, network: Object },
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

    const pwrDisplay = computed(() => (props.totalPower || 0).toFixed(0))
    const pwrPct = computed(() => Math.min(((props.totalPower || 0) / 500) * 100, 100))

    const netDown = computed(() => {
      const v = (props.network && props.network.download_mbps) || 0
      return v >= 100 ? v.toFixed(0) : v.toFixed(1)
    })
    const netUp = computed(() => {
      const v = (props.network && props.network.upload_mbps) || 0
      return v >= 100 ? v.toFixed(0) : v.toFixed(1)
    })
    const netPct = computed(() => {
      const val = Math.max(
        (props.network && props.network.download_mbps) || 0,
        (props.network && props.network.upload_mbps) || 0
      )
      return Math.min((val / 100) * 100, 100)
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

    return { topbar, memPct, pwrDisplay, pwrPct, netDown, netUp, netPct, fmt, barColor }
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
  padding: 16px 26px;
  flex-shrink: 0;
  cursor: grab;
  user-select: none;
  -webkit-user-select: none;
}

.clock { font-size: 36px; font-weight: 700; letter-spacing: 2px; color: var(--cyan); }
.date { font-size: 14px; color: var(--text-secondary); margin-top: 2px; }

.topbar-center { text-align: center; }
.uptime-label { font-size: 12px; color: var(--text-dim); letter-spacing: 2px; margin-bottom: 2px; }
.uptime-value { font-size: 18px; font-weight: 600; color: var(--text-primary); }
.sys-info { font-size: 11px; color: var(--text-dim); margin-top: 4px; }
.sys-sep { margin: 0 6px; opacity: 0.5; }
.sys-fps { color: var(--green); font-weight: 600; }

.topbar-right { display: flex; gap: 24px; }
.stat-item { text-align: center; }
.stat-label { font-size: 12px; color: var(--text-dim); display: block; letter-spacing: 1px; }
.stat-value { font-size: 20px; font-weight: 700; display: block; }

.mini-bar {
  width: 48px; height: 3px; background: var(--border); border-radius: 2px;
  margin-top: 4px; overflow: hidden;
}
.mini-bar-ghost { visibility: hidden; }
.mini-fill { height: 100%; border-radius: 2px; transition: width 0.3s; }
.pwr-val { color: var(--amber); font-size: 20px; white-space: nowrap; }
.pwr-unit { font-size: 13px; opacity: 0.7; }
.pwr-fill { background: linear-gradient(90deg, var(--amber), #fbbf24); }
.net-val { color: var(--blue); font-size: 16px; white-space: nowrap; }
.net-unit { font-size: 11px; opacity: 0.7; margin-left: 2px; }
.net-fill { background: linear-gradient(90deg, var(--blue), #93c5fd); }
.proc-val { color: var(--text-primary); font-size: 20px; }
</style>
