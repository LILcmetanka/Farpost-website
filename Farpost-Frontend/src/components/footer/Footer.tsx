function Footer() {
  return (
    <div className="flex flex-col flex-wrap sm:flex-row [@media(min-width:1000px)]:justify-between justify-start items-start mb-[30px] gap-15 mx-5 xl:mx-25">
      <div className="flex flex-col gap-[25px]">
        <h3 className="text-[16px]">Разделы</h3>
        <ul className="text-[12px] text-footer-text underline">
          <li><a href="#">Главная</a></li>
          <li><a href="#">Показания счетчиков</a></li>
          <li><a href="#">Уведомления об отключениях в мобильный</a></li>
          <li><a href="#">Управляющие организации Владивостока</a></li>
          <li><a href="#">Посты</a></li>
        </ul>
      </div>

      <div className="flex flex-col gap-[25px]">
        <h3 className="text-[16px]">Сайты Владивостока</h3>
        <ul className="text-[12px] text-footer-text underline">
          <li><a href="#">Авто</a></li>
          <li><a href="#">Кино</a></li>
          <li><a href="#">Недвижимость</a></li>
          <li><a href="#">Новости</a></li>
          <li><a href="#">Объявления</a></li>
          <li><a href="#">Работа</a></li>
          <li><a href="#">Справочник компаний</a></li>
        </ul>
      </div>

      <div className="flex flex-col gap-[25px]">
        <h3 className="text-[16px]">Информация</h3>
        <ul className="text-[12px] text-footer-text underline">
          <li><a href="#">Вакансии</a></li>
          <li><a href="#">Контакты</a></li>
        </ul>
        <div className="text-[10px] text-info-text">
          <p>По вопросам, предложениям иои ошибкам пишите на почту: <a className="text-black underline" href="#">off@vl.ru</a></p> <br />
          <p>© ООО «Фарпост», 2003 — 2025 <br />При любом использовании материалов <a className="text-black underline" href="#">ссылка на VL.ru</a> обязательна. <br />Цитирование в Интернете возможно только при наличии гиперссылки. <br />
          Все права защищены.</p>
        </div>
      </div>
    </div>
  );
}

export default Footer;