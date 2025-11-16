<template>
  <div class="animal-form-container">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? $t('animal.editAnimal') : $t('animal.addAnimal') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit">
      <Card>
        <template #title>
          {{ $t('animal.basicInfo') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="name">{{ $t('animal.name') }} *</label>
              <InputText
                id="name"
                v-model="formData.name"
                required
              />
            </div>

            <div class="form-field">
              <label for="name_en">{{ $t('animal.nameEn') }}</label>
              <InputText
                id="name_en"
                v-model="formData.name_en"
              />
            </div>

            <div class="form-field">
              <label for="species">{{ $t('animal.species') }} *</label>
              <Dropdown
                id="species"
                v-model="formData.species"
                :options="speciesOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="breed">{{ $t('animal.breed') }}</label>
              <InputText
                id="breed"
                v-model="formData.breed"
              />
            </div>

            <div class="form-field">
              <label for="category">{{ $t('animal.category') }} *</label>
              <Dropdown
                id="category"
                v-model="formData.category"
                :options="categoryOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="sex">{{ $t('animal.gender') }} *</label>
              <Dropdown
                id="sex"
                v-model="formData.sex"
                :options="sexOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="date_of_birth">{{ $t('animal.dateOfBirth') }}</label>
              <Calendar
                id="date_of_birth"
                v-model="formData.date_of_birth"
                date-format="yy-mm-dd"
              />
            </div>

            <div class="form-field">
              <label for="age_years">{{ $t('animal.ageYears') }}</label>
              <InputNumber
                id="age_years"
                v-model="formData.age_years"
                :min="0"
                :max="30"
              />
            </div>

            <div class="form-field">
              <label for="age_months">{{ $t('animal.ageMonths') }}</label>
              <InputNumber
                id="age_months"
                v-model="formData.age_months"
                :min="0"
                :max="11"
              />
            </div>

            <div class="form-field">
              <label for="color">{{ $t('animal.color') }}</label>
              <InputText
                id="color"
                v-model="formData.color"
              />
            </div>

            <div class="form-field">
              <label for="size">{{ $t('animal.size') }}</label>
              <Dropdown
                id="size"
                v-model="formData.size"
                :options="sizeOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field">
              <label for="weight">{{ $t('animal.weight') }}</label>
              <InputNumber
                id="weight"
                v-model="formData.weight"
                :min="0"
                :max-fraction-digits="2"
                suffix=" kg"
              />
            </div>

            <div class="form-field">
              <label for="microchip_id">{{ $t('animal.microchipId') }}</label>
              <InputText
                id="microchip_id"
                v-model="formData.microchip_id"
              />
            </div>

            <div class="form-field">
              <label for="status">{{ $t('animal.status') }} *</label>
              <Dropdown
                id="status"
                v-model="formData.status"
                :options="statusOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="intake_date">{{ $t('animal.intakeDate') }} *</label>
              <Calendar
                id="intake_date"
                v-model="formData.intake_date"
                date-format="yy-mm-dd"
                required
              />
            </div>

            <div class="form-field full-width">
              <label for="intake_reason">{{ $t('animal.intakeReason') }}</label>
              <Textarea
                id="intake_reason"
                v-model="formData.intake_reason"
                rows="3"
              />
            </div>
          </div>
        </template>
      </Card>

      <Card class="mt-3">
        <template #title>
          {{ $t('animal.medicalInfo') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label
                for="spayed_neutered"
                class="checkbox-label"
              >
                <Checkbox
                  id="spayed_neutered"
                  v-model="formData.spayed_neutered"
                  :binary="true"
                />
                {{ $t('animal.spayedNeutered') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="vaccinated"
                class="checkbox-label"
              >
                <Checkbox
                  id="vaccinated"
                  v-model="formData.vaccinated"
                  :binary="true"
                />
                {{ $t('animal.vaccinated') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="house_trained"
                class="checkbox-label"
              >
                <Checkbox
                  id="house_trained"
                  v-model="formData.house_trained"
                  :binary="true"
                />
                {{ $t('animal.houseTrained') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="good_with_kids"
                class="checkbox-label"
              >
                <Checkbox
                  id="good_with_kids"
                  v-model="formData.good_with_kids"
                  :binary="true"
                />
                {{ $t('animal.goodWithKids') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="good_with_dogs"
                class="checkbox-label"
              >
                <Checkbox
                  id="good_with_dogs"
                  v-model="formData.good_with_dogs"
                  :binary="true"
                />
                {{ $t('animal.goodWithDogs') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="good_with_cats"
                class="checkbox-label"
              >
                <Checkbox
                  id="good_with_cats"
                  v-model="formData.good_with_cats"
                  :binary="true"
                />
                {{ $t('animal.goodWithCats') }}
              </label>
            </div>

            <div class="form-field full-width">
              <label for="medical_history">{{ $t('animal.medicalHistory') }}</label>
              <Textarea
                id="medical_history"
                v-model="formData.medical_history"
                rows="4"
              />
            </div>

            <div class="form-field full-width">
              <label for="special_needs">{{ $t('animal.specialNeeds') }}</label>
              <Textarea
                id="special_needs"
                v-model="formData.special_needs"
                rows="3"
              />
            </div>
          </div>
        </template>
      </Card>

      <Card class="mt-3">
        <template #title>
          {{ $t('animal.description') }}
        </template>
        <template #content>
          <div class="form-field">
            <label for="description">{{ $t('animal.description') }} (PL)</label>
            <Textarea
              id="description"
              v-model="formData.description"
              rows="5"
            />
          </div>

          <div class="form-field">
            <label for="description_en">{{ $t('animal.descriptionEn') }}</label>
            <Textarea
              id="description_en"
              v-model="formData.description_en"
              rows="5"
            />
          </div>
        </template>
      </Card>

      <div class="form-actions">
        <Button
          type="button"
          :label="$t('common.cancel')"
          class="p-button-secondary"
          @click="router.back()"
        />
        <Button
          type="submit"
          :label="$t('common.save')"
          icon="pi pi-check"
          :loading="saving"
        />
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { animalService } from '@/services/animalService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'
import Checkbox from 'primevue/checkbox'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)
const saving = ref(false)

const formData = reactive({
  name: '',
  name_en: '',
  species: 'dog',
  breed: '',
  category: 'mammal',
  sex: 'unknown',
  date_of_birth: null,
  age_years: null,
  age_months: null,
  color: '',
  size: null,
  weight: null,
  microchip_id: '',
  status: 'available',
  intake_date: new Date(),
  intake_reason: '',
  description: '',
  description_en: '',
  medical_history: '',
  special_needs: '',
  spayed_neutered: false,
  vaccinated: false,
  house_trained: false,
  good_with_kids: false,
  good_with_dogs: false,
  good_with_cats: false
})

const speciesOptions = [
  { label: 'Dog', value: 'dog' },
  { label: 'Cat', value: 'cat' },
  { label: 'Rabbit', value: 'rabbit' },
  { label: 'Bird', value: 'bird' },
  { label: 'Other', value: 'other' }
]

const categoryOptions = [
  { label: 'Mammal', value: 'mammal' },
  { label: 'Bird', value: 'bird' },
  { label: 'Reptile', value: 'reptile' }
]

const sexOptions = ref([
  { label: t('animal.male'), value: 'male' },
  { label: t('animal.female'), value: 'female' },
  { label: t('animal.unknown'), value: 'unknown' }
])

const sizeOptions = ref([
  { label: t('animal.small'), value: 'small' },
  { label: t('animal.medium'), value: 'medium' },
  { label: t('animal.large'), value: 'large' },
  { label: t('animal.extraLarge'), value: 'extra_large' }
])

const statusOptions = ref([
  { label: t('animal.available'), value: 'available' },
  { label: t('animal.adopted'), value: 'adopted' },
  { label: t('animal.underTreatment'), value: 'under_treatment' },
  { label: t('animal.fostered'), value: 'fostered' },
  { label: t('animal.transferred'), value: 'transferred' }
])

const loadAnimal = async () => {
  if (!isEdit.value) return

  try {
    const animal = await animalService.getAnimal(route.params.id)
    Object.assign(formData, {
      ...animal,
      date_of_birth: animal.date_of_birth ? new Date(animal.date_of_birth) : null,
      intake_date: animal.intake_date ? new Date(animal.intake_date) : null
    })
  } catch (error) {
    console.error('Error loading animal:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load animal',
      life: 3000
    })
    router.push('/staff/animals')
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true

    const dataToSend = {
      ...formData,
      date_of_birth: formData.date_of_birth ? formData.date_of_birth.toISOString().split('T')[0] : null,
      intake_date: formData.intake_date ? formData.intake_date.toISOString().split('T')[0] : null
    }

    if (isEdit.value) {
      await animalService.updateAnimal(route.params.id, dataToSend)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: t('animal.animalUpdated'),
        life: 3000
      })
    } else {
      await animalService.createAnimal(dataToSend)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: t('animal.animalCreated'),
        life: 3000
      })
    }

    router.push('/staff/animals')
  } catch (error) {
    console.error('Error saving animal:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to save animal',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadAnimal()
})
</script>

<style scoped>
.animal-form-container {
  max-width: 1000px;
  margin: 0 auto;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.form-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field label {
  font-weight: 600;
  color: #374151;
}

.checkbox-label {
  flex-direction: row !important;
  align-items: center;
  gap: 0.75rem !important;
}

.full-width {
  grid-column: 1 / -1;
}

.mt-3 {
  margin-top: 1.5rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
