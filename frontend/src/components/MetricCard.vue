<template>
  <div class="metric-card">
    <div class="metric-row">
      <span class="metric-label">{{ label }}</span>
      <span class="metric-value">
        <span class="metric-num">{{ formattedValue }}</span>
        <span class="metric-unit">{{ unit }}</span>
      </span>
    </div>
    <div class="progress-bar">
      <div class="progress-fill" :style="fillStyle"></div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    label: { type: String, required: true },
    value: { type: Number, default: 0 },
    unit: { type: String, default: '' },
    max: { type: Number, default: 100 },
    barColor: { type: String, default: 'var(--cyan)' },
    decimals: { type: Number, default: 1 },
    showZero: { type: Boolean, default: true },
  },
  computed: {
    formattedValue() {
      if (!this.showZero && this.value === 0) return '--'
      return this.value.toFixed(this.decimals)
    },
    pct() {
      if (this.max <= 0) return 0
      const p = (this.value / this.max) * 100
      return Math.min(p, 100)
    },
    fillStyle() {
      return {
        width: this.pct + '%',
        background: this.barColor
      }
    }
  }
}
</script>

<style scoped>
.metric-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 14px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.metric-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 6px;
}

.metric-label {
  font-size: 12px;
  color: var(--text-dim);
  letter-spacing: 1px;
  white-space: nowrap;
}

.metric-value {
  display: flex;
  align-items: baseline;
  gap: 3px;
  white-space: nowrap;
}

.metric-num {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.metric-unit {
  font-size: 12px;
  color: var(--text-secondary);
}

.progress-bar {
  height: 3px;
  background: var(--border);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.4s ease;
}
</style>
