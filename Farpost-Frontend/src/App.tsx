import Calendar from "./components/calendar/Calendar"
import Charts from "./components/charts/Charts"
import JobsList from "./components/jobs-list/JobsList"
import ServiceOrgsList from "./components/service-orgs-list/ServiceOrgsList"
import SummaryToday from "./components/summary-today/SummaryToday"

function App() {
  return (
    <div className="flex flex-row gap-15">

      <div className="flex gap-[25px] flex-col">
        <SummaryToday />
        <Calendar />
      </div>

      <div className="flex flex-wrap gap-x-15 gap-y-[30px]">
        <ServiceOrgsList />
        <JobsList />
        <Charts />
      </div>
    </div>
  )
}

export default App
