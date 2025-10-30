import { useEffect, useState } from "react";
import NavbarButton from "./navbar-button/NavbarButton";
import ServiceOrgsItem from "./service-orgs-item/ServiceOrgsItem";
import type { Organization } from "../../types/organization.types";
import type { Service } from "../../types/service.types";

function ServiceOrgsList() {
  const [activeTab, setActiveTab] = useState<"services" | "orgs">("services");
  const [activeItem, setActiveItem] = useState<string | null>(null);
  const [orgs, setOrgs] = useState<Organization[]>([]);
  const [services, setServices] = useState<Service[]>([]);

  const [isLoading, setIsLoading] = useState(false);

  const currentList = activeTab === "services" ? services : orgs;

  useEffect(() => {
    const fetchData = async () => {
      const curr_time = new Date().toISOString();

      async function fetchOrganizations(curr_time?: string): Promise<Organization[]> {
        const baseUrl = import.meta.env.VITE_API_URL;
        let url = `${baseUrl}/off/orgs`;
        if (curr_time) {
          url += `?curr_time=${encodeURIComponent(curr_time)}`;
        }

        const response = await fetch(url);
        if (!response.ok) {
          throw new Error(`Ошибка запроса: ${response.status}`);
        }

        const data = await response.json();
        console.log("Raw organizations data:", data);

        return Array.isArray(data) ? data : data.organizations ?? [];
      }

      async function fetchServices(curr_time?: string): Promise<Service[]> {
        const baseUrl = import.meta.env.VITE_API_URL;
        let url = `${baseUrl}/off/blackouts`;
        if (curr_time) {
          url += `?curr_time=${encodeURIComponent(curr_time)}`;
        }

        const response = await fetch(url);
        if (!response.ok) {
          throw new Error(`Ошибка запроса: ${response.status}`);
        }

        const data = await response.json();
        console.log("Raw services data:", data);

        return Array.isArray(data) ? data : data.services ?? [];
      }

      try {
        setIsLoading(true);
        const [orgsData, servicesData] = await Promise.all([
          fetchOrganizations(curr_time),
          fetchServices(curr_time),
        ]);
        setOrgs(orgsData);
        setServices(servicesData);
      } catch (err) {
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);


  return (
    <div className="flex flex-col gap-[25px] max-w-[372px] w-full mx-5 sm:mx-0">
      <div className="flex gap-[25px]">
        <NavbarButton
          name="Коммунальные услуги"
          isActive={activeTab === "services"}
          onClick={() => setActiveTab("services")}
        />
        <NavbarButton
          name="Организации"
          isActive={activeTab === "orgs"}
          onClick={() => setActiveTab("orgs")}
        />
      </div>

      <div className="p-5 border border-line rounded-[10px] bg-black/1">
        {currentList.length === 0 ? (
          <div className="text-gray-500">Нет данных</div>
        ) : (
          <ul className="flex flex-col gap-5">
            {isLoading && <div>Загрзка данных</div>}
            {currentList.map((item) => (
              <ServiceOrgsItem
                key={"type" in item ? item.type : item.name}
                item={item}
                activeItem={activeItem}
                setActiveId={setActiveItem}
              />
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

export default ServiceOrgsList;
