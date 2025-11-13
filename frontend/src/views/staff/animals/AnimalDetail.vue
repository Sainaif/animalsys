<template>
  <div class="animal-detail">
    <LoadingSpinner v-if="loading" />

    <div v-else-if="animal" class="detail-container">
      <div class="detail-header">
        <Button
          icon="pi pi-arrow-left"
          class="p-button-text"
          @click="router.back()"
        />
        <h1>{{ animal.name }}</h1>
        <div class="header-actions">
          <Button
            :label="$t('animal.editAnimal')"
            icon="pi pi-pencil"
            @click="router.push(`/staff/animals/${animal.id}/edit`)"
          />
          <Button
            :label="$t('animal.deleteAnimal')"
            icon="pi pi-trash"
            class="p-button-danger"
            @click="confirmDelete"
          />
        </div>
      </div>

      <TabView>
        <TabPanel :header="$t('animal.basicInfo')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('animal.name') }}</label>
                  <p>{{ animal.name }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.species') }}</label>
                  <p>{{ animal.species }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.breed') }}</label>
                  <p>{{ animal.breed || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.gender') }}</label>
                  <p>{{ $t(`animal.${animal.sex}`) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.age') }}</label>
                  <p>{{ formatAge(animal) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.status') }}</label>
                  <Badge :variant="getStatusVariant(animal.status)">
                    {{ $t(`animal.${animal.status}`) }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.intakeDate') }}</label>
                  <p>{{ formatDate(animal.intake_date) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.color') }}</label>
                  <p>{{ animal.color || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.weight') }}</label>
                  <p>{{ animal.weight ? `${animal.weight} kg` : 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.microchipId') }}</label>
                  <p>{{ animal.microchip_id || 'N/A' }}</p>
                </div>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.description') }}</label>
                <p>{{ animal.description || 'N/A' }}</p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('animal.medicalInfo')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('animal.spayedNeutered') }}</label>
                  <Badge :variant="animal.spayed_neutered ? 'success' : 'neutral'">
                    {{ animal.spayed_neutered ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.vaccinated') }}</label>
                  <Badge :variant="animal.vaccinated ? 'success' : 'neutral'">
                    {{ animal.vaccinated ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.medicalHistory') }}</label>
                <p>{{ animal.medical_history || 'N/A' }}</p>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.specialNeeds') }}</label>
                <p>{{ animal.special_needs || 'N/A' }}</p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('animal.behaviorInfo')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('animal.goodWithKids') }}</label>
                  <Badge :variant="animal.good_with_kids ? 'success' : 'neutral'">
                    {{ animal.good_with_kids ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.goodWithDogs') }}</label>
                  <Badge :variant="animal.good_with_dogs ? 'success' : 'neutral'">
                    {{ animal.good_with_dogs ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.goodWithCats') }}</label>
                  <Badge :variant="animal.good_with_cats ? 'success' : 'neutral'">
                    {{ animal.good_with_cats ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.houseTrained') }}</label>
                  <Badge :variant="animal.house_trained ? 'success' : 'neutral'">
                    {{ animal.house_trained ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.temperament') }}</label>
                <div class="temperament-badges">
                  <Badge v-for="temp in animal.temperament" :key="temp" variant="primary">
                    {{ temp }}
                  </Badge>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>
      </TabView>
    </div>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { animalService } from '@/services/animalService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const animal = ref(null)
const loading = ref(true)

const loadAnimal = async () => {
  try {
    loading.value = true
    animal.value = await animalService.getAnimal(route.params.id)
  } catch (error) {
    console.error('Error loading animal:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load animal',
      life: 3000
    })
    router.push('/staff/animals')
  } finally {
    loading.value = false
  }
}

const formatAge = (animal) => {
  if (animal.age_years) return `${animal.age_years} ${t('animal.ageYears')}`
  if (animal.age_months) return `${animal.age_months} ${t('animal.ageMonths')}`
  return 'N/A'
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString()
}

const getStatusVariant = (status) => {
  const variants = {
    available: 'success',
    adopted: 'info',
    under_treatment: 'warning',
    fostered: 'primary',
    transferred: 'neutral',
    deceased: 'danger'
  }
  return variants[status] || 'neutral'
}

const confirmDelete = () => {
  confirm.require({
    message: t('animal.confirmDelete'),
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await animalService.deleteAnimal(animal.value.id)
        toast.add({
          severity: 'success',
          summary: 'Success',
          detail: t('animal.animalDeleted'),
          life: 3000
        })
        router.push('/staff/animals')
      } catch (error) {
        console.error('Error deleting animal:', error)
        toast.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to delete animal',
          life: 3000
        })
      }
    }
  })
}

onMounted(() => {
  loadAnimal()
})
</script>

<style scoped>
.animal-detail {
  max-width: 1200px;
  margin: 0 auto;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.detail-header h1 {
  flex: 1;
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item label {
  font-weight: 600;
  color: #6b7280;
  font-size: 0.875rem;
  text-transform: uppercase;
}

.info-item p {
  color: #2c3e50;
  font-size: 1rem;
  margin: 0;
}

.full-width {
  grid-column: 1 / -1;
}

.temperament-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
</style>
