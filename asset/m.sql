--select count(*) from available_view;
--select count(*) from medis;
--select * from available_view where "個別医薬品コード" like '3961008F1020';
select * from yj;
select "個別医薬品コード"='', "個別医薬品コード", custom_yj, yj_status  from available_view where "HOT11" like '18200030401';