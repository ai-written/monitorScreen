<template>
  <div class="block fan-block">
    <div class="block-header">
      <div class="block-title"><span class="dot fan-dot"></span> 风扇</div>
    </div>
    <div class="fan-list">
      <div v-if="!fans || fans.length === 0" class="fan-empty">无风扇数据</div>
      <div v-for="(f, i) in fans" :key="i" class="fan-item">
        <div class="fan-row">
          <span class="fan-name">{{ f.name || '风扇 ' + (i+1) }}</span>
          <span class="fan-stats">
            <span class="fan-rpm">{{ fmt(f.rpm) }} RPM</span>
            <span class="fan-pct" v-if="f.percent > 0">{{ f.percent.toFixed(0) }}%</span>
          </span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill fan-fill" :style="{ width: fanPct(f) + '%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: { fans: Array },
  methods: {
    fmt(v) { return typeof v === 'number' ? v.toFixed(0) : '--' },
    fanPct(f) { return f.percent > 0 ? f.percent : Math.min((f.rpm / 2000) * 100, 100) }
  }
}
</script>

<style scoped>
.block {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
}

.fan-block { border-top: 2px solid var(--pink); }

.block-header { margin-bottom: 10px; }
.block-title {
  font-size: 14px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.fan-dot { background: var(--pink); }

.fan-list { display: flex; flex-direction: column; gap: 10px; flex: 1; }
.fan-empty { color: var(--text-dim); font-size: 14px; text-align: center; padding: 20px 0; }

.fan-row { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 4px; }
.fan-name {
  font-size: 14px; font-weight: 600; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  flex: 1; min-width: 0;
}
.fan-stats { display: flex; gap: 6px; white-space: nowrap; flex-shrink: 0; }
.fan-rpm { font-size: 13px; color: var(--pink); font-weight: 600; }
.fan-pct { font-size: 13px; color: var(--text-secondary); }

.progress-bar {
  height: 6px; background: var(--border); border-radius: 3px; overflow: hidden;
}
.progress-fill { height: 100%; border-radius: 3px; transition: width 0.5s; }
.fan-fill { background: linear-gradient(90deg, var(--pink), #f9a8d4); }
</style>
