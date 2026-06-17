<template>
  <div class="block mem-block">
    <div class="block-header">
      <div class="block-title"><span class="dot mem-dot"></span> MEMORY</div>
    </div>
    <div v-if="memory.brand" class="mem-brand">{{ memory.brand }}</div>
    <div class="mem-main">
      <div class="mem-used">{{ fmt(memory.used_gb) }}<span class="mem-unit"> GB</span></div>
      <div class="mem-total">/ {{ memory.total_gb ? memory.total_gb.toFixed(1) : '--' }} GB</div>
    </div>
    <div class="progress-bar">
      <div class="progress-fill mem-fill" :style="{ width: memPct + '%' }"></div>
    </div>
    <div class="mem-specs">
      <div class="spec" v-if="memory.type">{{ memory.type }}</div>
      <div class="spec" v-if="memory.frequency">{{ memory.frequency }}</div>
      <div class="spec" v-if="memory.channel">{{ memory.channel }}</div>
    </div>
  </div>
</template>

<script>
export default {
  props: { memory: Object },
  computed: {
    memPct() {
      if (this.memory && this.memory.total_gb > 0) {
        return (this.memory.used_gb / this.memory.total_gb) * 100
      }
      return 0
    }
  },
  methods: {
    fmt(v) { return typeof v === 'number' ? v.toFixed(1) : '--' }
  }
}
</script>

<style scoped>
.block {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
}

.mem-block { border-top: 2px solid var(--green); }

.block-header { margin-bottom: 4px; }
.block-title {
  font-size: 12px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.mem-dot { background: var(--green); }

.mem-brand {
  font-size: 11px; color: var(--text-dim);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 8px;
}

.mem-main { display: flex; align-items: baseline; gap: 6px; margin-bottom: 10px; }
.mem-used { font-size: 30px; font-weight: 700; color: var(--green); }
.mem-unit { font-size: 18px; color: var(--green); }
.mem-total { font-size: 16px; color: var(--text-secondary); }

.progress-bar {
  height: 8px; background: var(--border); border-radius: 4px;
  overflow: hidden; margin-bottom: 12px;
}
.progress-fill { height: 100%; border-radius: 4px; transition: width 0.5s; }
.mem-fill { background: linear-gradient(90deg, var(--green), #22c55e); }

.mem-specs { display: flex; gap: 8px; margin-top: auto; flex-wrap: wrap; }
.spec {
  font-size: 12px; color: var(--text-secondary);
  background: var(--green-dim); padding: 3px 10px; border-radius: 6px;
}
</style>
